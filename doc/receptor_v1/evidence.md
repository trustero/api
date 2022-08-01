# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [receptor_v1/evidence.proto](#receptor_v1_evidence-proto)
    - [Evidence](#receptor_v1-Evidence)
    - [Evidence.Source](#receptor_v1-Evidence-Source)
  
- [Scalar Value Types](#scalar-value-types)



<a name="receptor_v1_evidence-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## receptor_v1/evidence.proto
This file is subject to the terms and conditions defined in
file &#39;LICENSE.txt&#39;, which is part of this source code package.


<a name="receptor_v1-Evidence"></a>

### Evidence
An evidence is a unstructured or structured document that represent the how a service is being used within a
service provider account.  For example, the configuration of an S3 bucket in AWS.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| caption | [string](#string) |  | Human readable English string that identifies this evidence. It&#39;s important the caption is stable for all scans of the same evidence type. |
| description | [string](#string) |  | Human readable English string describing the content of this evidence. |
| service_name | [string](#string) |  | The name of service this evidence was collected from. For example, &#34;S3&#34;. The service name must be one of the service types reported in Services struct (See the message Service definition). |
| sources | [Evidence.Source](#receptor_v1-Evidence-Source) | repeated |  |
| doc | [Document](#receptor_v1-Document) |  | An unstructured evidence. |
| struct | [Struct](#receptor_v1-Struct) |  | A structured evidence |






<a name="receptor_v1-Evidence-Source"></a>

### Evidence.Source
The raw API request used to generate this evidence.  The raw API request and response are used to prove to
examiners this evidence correlates to real service instance configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| raw_api_request | [string](#string) |  |  |
| raw_api_response | [string](#string) |  | The raw API response used to generate this evidence. The raw API request and response is used to prove to examiners this evidence correlates to real service instance configuration. |





 

 

 

 



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

