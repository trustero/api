# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [receptor_v1/finding.proto](#receptor_v1_finding-proto)
    - [Finding](#receptor_v1-Finding)
  
- [Scalar Value Types](#scalar-value-types)



<a name="receptor_v1_finding-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## receptor_v1/finding.proto
This file is subject to the terms and conditions defined in
file &#39;LICENSE.txt&#39;, which is part of this source code package.


<a name="receptor_v1-Finding"></a>

### Finding
A finding is a set of evidence(s) collected from a service provider account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_type | [string](#string) |  | Unique receptor identifier. A receptor is expected to report findings from only one service provider type. A stable identifier that represent the type of receptor reporting this finding. The identifier is akin to a fully qualified Go package name or a Java class name. For example, &#34;github.com/trustero/receptor/gitlab&#34;. REMIND maps to Receptor.ModelID |
| service_provider_account | [string](#string) |  | The receptor&#39;s evidence source. REMIND maps to Receptor.TenantID |
| evidences | [Evidence](#receptor_v1-Evidence) | repeated | One or more evidence collected by a typical receptor scan. |





 

 

 

 



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

