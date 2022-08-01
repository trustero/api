# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [receptor_v1/struct.proto](#receptor_v1_struct-proto)
    - [Struct](#receptor_v1-Struct)
    - [Struct.ColDisplayNamesEntry](#receptor_v1-Struct-ColDisplayNamesEntry)
    - [Struct.Row](#receptor_v1-Struct-Row)
    - [Struct.Row.ColsEntry](#receptor_v1-Struct-Row-ColsEntry)
    - [Struct.Row.Value](#receptor_v1-Struct-Row-Value)
  
- [Scalar Value Types](#scalar-value-types)



<a name="receptor_v1_struct-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## receptor_v1/struct.proto
This file is subject to the terms and conditions defined in
file &#39;LICENSE.txt&#39;, which is part of this source code package.


<a name="receptor_v1-Struct"></a>

### Struct
A structured evidence defined in tabular form.  Each struct typically represent a service type (see Evidence
message definition).  Each struct consists of rows of data.  Each row typically represent a service instance and
its configurations.  A row  contains column name and column value pairs.  All rows in a struct must have the same
column name-value pairs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rows | [Struct.Row](#receptor_v1-Struct-Row) | repeated | Each row typically represents the configuration of a service instance. @required |
| col_display_names | [Struct.ColDisplayNamesEntry](#receptor_v1-Struct-ColDisplayNamesEntry) | repeated | A map of row column name to display name pairs. @required |
| col_display_order | [string](#string) | repeated | An ordered list of row column names. The order defines how each column will be rendered by default. @required |






<a name="receptor_v1-Struct-ColDisplayNamesEntry"></a>

### Struct.ColDisplayNamesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="receptor_v1-Struct-Row"></a>

### Struct.Row
A row of structured data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [string](#string) |  | A unique service ID relative to the struct. A row ID typically represents a unique service ID. The service_id must be previously reported in the Services message. @required |
| cols | [Struct.Row.ColsEntry](#receptor_v1-Struct-Row-ColsEntry) | repeated | Columns of the row in column name to value pairs. All rows in a struct must have the same column names and corresponding value types. @required |






<a name="receptor_v1-Struct-Row-ColsEntry"></a>

### Struct.Row.ColsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Struct.Row.Value](#receptor_v1-Struct-Row-Value) |  |  |






<a name="receptor_v1-Struct-Row-Value"></a>

### Struct.Row.Value
Column value types can be any protobuf scalar or google.proto.Timestamp.
@required


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| double_value | [double](#double) |  |  |
| float_value | [float](#float) |  |  |
| int32_value | [int32](#int32) |  |  |
| int64_value | [int64](#int64) |  |  |
| uint32_value | [uint32](#uint32) |  |  |
| uint64_value | [uint64](#uint64) |  |  |
| sint32_value | [sint32](#sint32) |  |  |
| sint64_value | [sint64](#sint64) |  |  |
| fixed32_value | [fixed32](#fixed32) |  |  |
| fixed64_value | [fixed64](#fixed64) |  |  |
| sfixed32_value | [sfixed32](#sfixed32) |  |  |
| sfixed64_value | [sfixed64](#sfixed64) |  |  |
| bool_value | [bool](#bool) |  |  |
| string_value | [string](#string) |  |  |
| timestamp_value | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

