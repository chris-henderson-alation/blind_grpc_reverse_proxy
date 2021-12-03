# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: test.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='test.proto',
  package='test',
  syntax='proto3',
  serialized_options=b'Z\006.;main',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\ntest.proto\x12\x04test\"\x19\n\x06String\x12\x0f\n\x07message\x18\x01 \x01(\t2\xc0\x01\n\x04Test\x12%\n\x05Unary\x12\x0c.test.String\x1a\x0c.test.String\"\x00\x12.\n\x0cServerStream\x12\x0c.test.String\x1a\x0c.test.String\"\x00\x30\x01\x12.\n\x0c\x43lientStream\x12\x0c.test.String\x1a\x0c.test.String\"\x00(\x01\x12\x31\n\rBidirectional\x12\x0c.test.String\x1a\x0c.test.String\"\x00(\x01\x30\x01\x42\x08Z\x06.;mainb\x06proto3'
)




_STRING = _descriptor.Descriptor(
  name='String',
  full_name='test.String',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='message', full_name='test.String.message', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=20,
  serialized_end=45,
)

DESCRIPTOR.message_types_by_name['String'] = _STRING
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

String = _reflection.GeneratedProtocolMessageType('String', (_message.Message,), {
  'DESCRIPTOR' : _STRING,
  '__module__' : 'test_pb2'
  # @@protoc_insertion_point(class_scope:test.String)
  })
_sym_db.RegisterMessage(String)


DESCRIPTOR._options = None

_TEST = _descriptor.ServiceDescriptor(
  name='Test',
  full_name='test.Test',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=48,
  serialized_end=240,
  methods=[
  _descriptor.MethodDescriptor(
    name='Unary',
    full_name='test.Test.Unary',
    index=0,
    containing_service=None,
    input_type=_STRING,
    output_type=_STRING,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='ServerStream',
    full_name='test.Test.ServerStream',
    index=1,
    containing_service=None,
    input_type=_STRING,
    output_type=_STRING,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='ClientStream',
    full_name='test.Test.ClientStream',
    index=2,
    containing_service=None,
    input_type=_STRING,
    output_type=_STRING,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Bidirectional',
    full_name='test.Test.Bidirectional',
    index=3,
    containing_service=None,
    input_type=_STRING,
    output_type=_STRING,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_TEST)

DESCRIPTOR.services_by_name['Test'] = _TEST

# @@protoc_insertion_point(module_scope)
