package main

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

func BsonRegistry() *bsoncodec.Registry {
	registry := bson.NewRegistry()
	registerUUID(registry)

	return registry
}

var errBadUUID = errors.New("incorrect type, UUID expected")

func registerUUID(registry *bsoncodec.Registry) {
	const mongoUUIDBinarySubtype = byte(0x04)

	registry.RegisterTypeEncoder(
		reflect.TypeOf(uuid.UUID{}),
		bsoncodec.ValueEncoderFunc(
			func(_ bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
				if b, ok := value.Interface().(uuid.UUID); ok {
					return errors.Wrap(writer.WriteBinaryWithSubtype(b[:], mongoUUIDBinarySubtype), "write uuid")
				}

				return errBadUUID
			}),
	)

	registry.RegisterTypeDecoder(
		reflect.TypeOf(uuid.UUID{}),
		bsoncodec.ValueDecoderFunc(
			func(_ bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
				data, _, err := reader.ReadBinary()
				if err != nil {
					return errors.Wrap(err, "read uuid")
				}

				val, err := uuid.FromBytes(data)
				if err != nil {
					return errors.Wrap(err, "decode uuid")
				}

				value.Set(reflect.ValueOf(val))

				return nil
			}),
	)
}
