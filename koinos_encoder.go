package log

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	_pool = buffer.NewPool()
	// GetPool retrieves a buffer from the pool, creating one if necessary.
	GetPool = _pool.Get
)

// ----------------------------------------------------------------------------
// sliceArrayEncoder clone
// ----------------------------------------------------------------------------

// sliceArrayEncoder is an ArrayEncoder backed by a simple []interface{}. Like
// the MapObjectEncoder, it's not designed for production use.
type sliceArrayEncoder struct {
	elems []interface{}
}

func (s *sliceArrayEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	enc := &sliceArrayEncoder{}
	err := v.MarshalLogArray(enc)
	s.elems = append(s.elems, enc.elems)
	return err
}

func (s *sliceArrayEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	m := zapcore.NewMapObjectEncoder()
	err := v.MarshalLogObject(m)
	s.elems = append(s.elems, m.Fields)
	return err
}

func (s *sliceArrayEncoder) AppendReflected(v interface{}) error {
	s.elems = append(s.elems, v)
	return nil
}

func (s *sliceArrayEncoder) AppendBool(v bool)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendByteString(v []byte)      { s.elems = append(s.elems, string(v)) }
func (s *sliceArrayEncoder) AppendComplex128(v complex128)  { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendComplex64(v complex64)    { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendDuration(v time.Duration) { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendFloat64(v float64)        { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendFloat32(v float32)        { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt(v int)                { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt64(v int64)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt32(v int32)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt16(v int16)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt8(v int8)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendString(v string)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendTime(v time.Time)         { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint(v uint)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint64(v uint64)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint32(v uint32)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint16(v uint16)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint8(v uint8)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUintptr(v uintptr)        { s.elems = append(s.elems, v) }

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

// KoinosEncoder implements custon koinos log formatting
type KoinosEncoder struct {
	*zapcore.MapObjectEncoder
	*zapcore.EncoderConfig
	ConsoleSeparator string
	AppID            string
}

// NewKoinosEncoder creates and returns a new instance of KoinosEncoder
func NewKoinosEncoder(cfg zapcore.EncoderConfig, appID string) zapcore.Encoder {
	return &KoinosEncoder{EncoderConfig: &cfg, ConsoleSeparator: "", AppID: appID}
}

// Clone clones KoinosEncoder
func (ke *KoinosEncoder) Clone() zapcore.Encoder {
	return &KoinosEncoder{}
}

// EncodeEntry encodes the given entry data
func (ke *KoinosEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := GetPool()

	arr := getSliceEncoder()
	if ke.TimeKey != "" && ke.EncodeTime != nil {
		ke.EncodeTime(ent.Time, arr)
		arr.AppendString(" ")
	}

	arr.AppendString("(")
	arr.AppendString(ke.AppID)
	arr.AppendString(")")

	arr.AppendString(" [")
	if ent.Caller.Defined {
		if ke.CallerKey != "" && ke.EncodeCaller != nil {
			ke.EncodeCaller(ent.Caller, arr)
		}
	}

	arr.AppendString("] <")
	if ke.LevelKey != "" && ke.EncodeLevel != nil {
		ke.EncodeLevel(ent.Level, arr)
	}

	arr.AppendString(">: ")

	for i := range arr.elems {
		if i > 0 {
			line.AppendString(ke.ConsoleSeparator)
		}
		fmt.Fprint(line, arr.elems[i])
	}
	putSliceEncoder(arr)

	if ke.MessageKey != "" {
		line.AppendString(ent.Message)
	}

	if ke.LineEnding != "" {
		line.AppendString(ke.LineEnding)
	} else {
		line.AppendString("\n")
	}

	return line, nil
}

const (
	red color = iota + 31
	green
	yellow
	blue
)

type color uint8

var (
	_levelToColor = map[zapcore.Level]color{
		zapcore.DebugLevel:  blue,
		zapcore.InfoLevel:   green,
		zapcore.WarnLevel:   yellow,
		zapcore.ErrorLevel:  red,
		zapcore.DPanicLevel: red,
		zapcore.PanicLevel:  red,
		zapcore.FatalLevel:  red,
	}
	_unknownLevelColor = red

	_koinosColorString = make(map[zapcore.Level]string, len(_levelToColor))
)

func (c color) AddColor(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

func init() {
	for level, color := range _levelToColor {
		if level == zapcore.WarnLevel {
			_koinosColorString[level] = color.AddColor("warning")
		} else {
			_koinosColorString[level] = color.AddColor(level.String())
		}
	}
}

// KoinosLevelEncoder implements the Koinos log level encoding standard
func KoinosLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == zapcore.WarnLevel {
		enc.AppendString("warning")
	} else {
		enc.AppendString(l.String())
	}
}

// KoinosColorLevelEncoder implements the Koinos log level color encoding standard
func KoinosColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _koinosColorString[l]
	if !ok {
		s = _unknownLevelColor.AddColor(l.String())
	}
	enc.AppendString(s)
}

// KoinosTimeEncoder encodes timestamps in the Koinos human-readable log standard
func KoinosTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}
