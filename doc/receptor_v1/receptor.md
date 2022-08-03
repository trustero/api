# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [receptor_v1/receptor.proto](#receptor_v1_receptor-proto)
    - [Credential](#receptor_v1-Credential)
    - [Document](#receptor_v1-Document)
    - [Evidence](#receptor_v1-Evidence)
    - [Evidence.Source](#receptor_v1-Evidence-Source)
    - [Finding](#receptor_v1-Finding)
    - [JobResult](#receptor_v1-JobResult)
    - [ReceptorConfiguration](#receptor_v1-ReceptorConfiguration)
    - [ReceptorOID](#receptor_v1-ReceptorOID)
    - [Services](#receptor_v1-Services)
    - [Services.Service](#receptor_v1-Services-Service)
    - [Struct](#receptor_v1-Struct)
    - [Struct.ColDisplayNamesEntry](#receptor_v1-Struct-ColDisplayNamesEntry)
    - [Struct.Row](#receptor_v1-Struct-Row)
    - [Struct.Row.ColsEntry](#receptor_v1-Struct-Row-ColsEntry)
    - [Struct.Row.Value](#receptor_v1-Struct-Row-Value)
  
    - [Receptor](#receptor_v1-Receptor)
  
- [Scalar Value Types](#scalar-value-types)



<a name="receptor_v1_receptor-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## receptor_v1/receptor.proto
This file is subject to the terms and conditions defined in
file &#39;LICENSE.txt&#39;, which is part of this source code package.


<a name="receptor_v1-Credential"></a>

### Credential
Credential to access a service provider account.
REMIND:  Credential maps to receptor.VerifyResult record with the addition of credential being verified.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_object_id | [string](#string) |  | Trustero&#39;s receptor record identifier. This identifier is typically provided to the receptor as part of a reporting findings or discover services request. |
| credential | [string](#string) |  | The service provider credential being verified. |
| is_credential_valid | [bool](#bool) |  | Report whether the service provider credential provided in this message is valid for report findings or discover services request. |
| message | [string](#string) |  | Reason for why the service provider credential in this message is invalid. |






<a name="receptor_v1-Document"></a>

### Document
An unstructured evidence provided as a MIME document.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [string](#string) |  | A unique service ID relative to the document. A row ID typically represents a unique service ID. The service_id must be previously reported in the Services message. @required |
| mime | [string](#string) |  | Document type defined using MIME (https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types). @required |
| body | [bytes](#bytes) |  | Opaque document body. The document body must match the type defined by the mime attribute. @required |






<a name="receptor_v1-Evidence"></a>

### Evidence
An evidence is a unstructured or structured document that represent the how a service is being used within a
service provider account.  For example, the configuration of an S3 bucket in AWS.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| caption | [string](#string) |  | Human readable English string that identifies this evidence. It&#39;s important the caption is stable for all scans of the same evidence type. |
| description | [string](#string) |  | Human readable English string describing the content of this evidence. |
| service_name | [string](#string) |  | The name of service this evidence was collected from. For example, &#34;S3&#34;. The service name must be one of the service types reported in Services struct (See the message Service definition). |
| sources | [Evidence.Source](#receptor_v1-Evidence-Source) | repeated | The raw API request used to generate this evidence. The raw API request and response are used to prove to examiners this evidence correlates to real service instance configuration. |
| doc | [Document](#receptor_v1-Document) |  | An unstructured evidence. |
| struct | [Struct](#receptor_v1-Struct) |  | A structured evidence |






<a name="receptor_v1-Evidence-Source"></a>

### Evidence.Source



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| raw_api_request | [string](#string) |  | The raw API request used to generate this evidence. The raw API request and response are used to prove to examiners this evidence correlates to real service instance configuration. |
| raw_api_response | [string](#string) |  | The raw API response used to generate this evidence. The raw API request and response is used to prove to examiners this evidence correlates to real service instance configuration. |






<a name="receptor_v1-Finding"></a>

### Finding
A finding is a set of evidence(s) collected from a service provider account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_type | [string](#string) |  | Unique receptor identifier. A receptor is expected to report findings from only one service provider type. A stable identifier that represent the type of receptor reporting this finding. The identifier is a simple URL encoded string that includes an organization name and the service provider name. For example: &#34;trustero_gitlab&#34;. |
| service_provider_account | [string](#string) |  | The receptor&#39;s evidence source. REMIND maps to Receptor.TenantID |
| evidences | [Evidence](#receptor_v1-Evidence) | repeated | One or more evidence collected by a typical receptor scan. |






<a name="receptor_v1-JobResult"></a>

### JobResult
Trustero uses asynchronous jobs to track receptor requests.  Trustero initiates a receptor job providing a
receptor_object_id, a tracer_id, and a command.  When the receptor completes the job, the receptor callback
into Trustero to report the job result.
REMIND:  JobResult maps to AsyncTask


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tracer_id | [string](#string) |  | A tracer ID used to track the progress of the receptor request. REMIND AyncTask.TracerID for tracking. |
| command | [string](#string) |  | Receptor command request that completed. One of &#34;verify&#34;, &#34;scan&#34;, or &#34;discover&#34; |
| result | [string](#string) |  | Receptor command request result. One of &#34;success&#34;, &#34;fail&#34;, or &#34;error&#34;. |
| receptor_object_id | [string](#string) |  | Trustero&#39;s receptor record identifier. REMIND Receptor.ID |






<a name="receptor_v1-ReceptorConfiguration"></a>

### ReceptorConfiguration
Trustero stored receptor configuration and service provider credential.
REMIND: ReceptorConfiguration is a subset of existing ntrced&#39;s Receptor record.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_object_id | [string](#string) |  | Trustero receptor record identifier. REMIND Receptor.ID |
| credential | [string](#string) |  | Credential required to access a service provider for report finding and discover services purposes. REMIND Receptor.Credential required to access the target service. |
| config | [string](#string) |  | Additional receptor configuration to access a service provider account. REMIND Receptor.config task configuration in json. |
| service_provider_account | [string](#string) |  | Service provider account REMIND Receptor.TenantID |






<a name="receptor_v1-ReceptorOID"></a>

### ReceptorOID
Trustero receptor record identifier.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_object_id | [string](#string) |  | Trustero string representation of a persistent record. |






<a name="receptor_v1-Services"></a>

### Services
Service instances configured within a service provider account.  For example, all service instances configured in
an AWS account which may include S3 buckets, ECS clusters, RDS database instances, etc.  The boundary of a
service instance such as a ECS cluster or an ECS container instance is dependent on how the findings are
collected.  Each service instance_id should be associated with at least one Evidence.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_type | [string](#string) |  | Unique receptor type. A stable string identifier that represent the type of receptor reporting this finding. The identifier is a simple URL encode string that includes the organization name and a service provider name. For example &#34;trustero_gitlab&#34;. @required |
| service_provider_account | [string](#string) |  | The service provider of this list of services. @required |
| services | [Services.Service](#receptor_v1-Services-Service) | repeated | A list of service instances. @required |






<a name="receptor_v1-Services-Service"></a>

### Services.Service
A service instance definition.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the service. For example, &#34;ECS&#34;. @required |
| instance_id | [string](#string) |  | Unique service ID. For example, ECS&#39;s UUID. @required |






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
| bool_value | [bool](#bool) |  |  |
| string_value | [string](#string) |  |  |
| timestamp_value | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 

 

 


<a name="receptor_v1-Receptor"></a>

### Receptor
A receptor, or a Trustero client application, collects findings supporting the use of a service from a service
provider instance.  An example of a service provider is AWS and an example of a service provider account is an
AWS account.  An example of a service is S3 and an example of a service instance is an S3 bucket.  Trustero
associates the collected evidence to support the fact an organization is following it&#39;s stated practices.  A finding
is comprised of a list of evidences.  Each evidence is associated with a service instance and contains its
configuration information. An example of a finding is an AWS S3 bucket and its configuration.  Service configuration
can be in opaque document format or structured document format.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Verified | [Credential](#receptor_v1-Credential) | [.google.protobuf.Empty](#google-protobuf-Empty) | Report whether the provided credential is a valid service provider credential for purpose of discovering services and reporting findings. This rpc call is typically made as callback by a receptor to trustero from a check credential receptor request. |
| GetConfiguration | [ReceptorOID](#receptor_v1-ReceptorOID) | [ReceptorConfiguration](#receptor_v1-ReceptorConfiguration) | Get the receptor configuration and service provider credential using the provided receptor record identifier. This rpc call is typically made as a callback by a receptor prior to making a report findings or discover services receptor request. |
| Discovered | [Services](#receptor_v1-Services) | [.google.protobuf.StringValue](#google-protobuf-StringValue) | Report known services. A receptor or a Trustero client application reports its known services on demand. This call returns a string value service listing ID or an error specifying why Trustero failed to process the service listing. |
| Report | [Finding](#receptor_v1-Finding) | [.google.protobuf.StringValue](#google-protobuf-StringValue) | Report a finding. A receptor or a Trustero client application reports its findings on a periodic basis. This call returns a string value collection ID or an error specifying why Trustero failed to process the finding. |
| Notify | [JobResult](#receptor_v1-JobResult) | [.google.protobuf.Empty](#google-protobuf-Empty) | Notify Trustero a long running report finding or discover services receptor request has completed. JobResult contains information about the receptor request and it&#39;s corresponding result. Information such as the JobResult.receptor_object_id are passed to the receptor as part of the request. |

 



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

