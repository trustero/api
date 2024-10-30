# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [receptor_v1/receptor.proto](#receptor_v1_receptor-proto)
    - [Credential](#receptor_v1-Credential)
    - [Document](#receptor_v1-Document)
    - [Evidence](#receptor_v1-Evidence)
    - [Finding](#receptor_v1-Finding)
    - [JobResult](#receptor_v1-JobResult)
    - [ReceptorConfiguration](#receptor_v1-ReceptorConfiguration)
    - [ReceptorOID](#receptor_v1-ReceptorOID)
    - [ReportChunk](#receptor_v1-ReportChunk)
    - [ReportResponse](#receptor_v1-ReportResponse)
    - [Row](#receptor_v1-Row)
    - [Row.ColsEntry](#receptor_v1-Row-ColsEntry)
    - [ServiceEntities](#receptor_v1-ServiceEntities)
    - [ServiceEntity](#receptor_v1-ServiceEntity)
    - [Source](#receptor_v1-Source)
    - [Struct](#receptor_v1-Struct)
    - [Struct.ColDisplayNamesEntry](#receptor_v1-Struct-ColDisplayNamesEntry)
    - [Struct.ColTagsEntry](#receptor_v1-Struct-ColTagsEntry)
    - [Value](#receptor_v1-Value)
  
    - [Receptor](#receptor_v1-Receptor)
  
- [Scalar Value Types](#scalar-value-types)



<a name="receptor_v1_receptor-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## receptor_v1/receptor.proto
This file is subject to the terms and conditions defined in
file &#39;LICENSE.txt&#39;, which is part of this source code package.


<a name="receptor_v1-Credential"></a>

### Credential
Credential is returned by a Verified request noting if a given service provider account credential is  valid.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_object_id | [string](#string) |  | Receptor_object_id is Trustero&#39;s receptor record identifier. This identifier is typically provided to the receptor as part of a reporting findings or discover services request. |
| credential | [string](#string) |  | Credential is the service provider credential being verified. |
| is_credential_valid | [bool](#bool) |  | Is_credential_valid report whether the service provider credential provided to the receptor verify request is valid. |
| message | [string](#string) |  | Message contains the reason for why the service provider credential in this message is invalid. |
| exceptions | [string](#string) |  | Exceptions contains information about the permissions that are missing for the credentials provided. |






<a name="receptor_v1-Document"></a>

### Document
Document is an unstructured evidence provided as a MIME document.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mime | [string](#string) |  | Mime is the document type defined using [MIME]. [MIME]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types |
| body | [bytes](#bytes) |  | Body is the opaque document body. The document body must match the type defined by the mime attribute. |






<a name="receptor_v1-Evidence"></a>

### Evidence
Evidence is a unstructured or structured document that represent the how a service is being used within a
service provider account.  For example, the configuration of an S3 bucket in AWS.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| caption | [string](#string) |  | Caption is a human readable English string that identifies this evidence. Caption must be stable for all scans of the same evidence type. Trustero uses the caption to associate this evidence with a set of relevant controls. |
| description | [string](#string) |  | Description is a human readable English string describing the content of this evidence. Description tells Trustero and users contents of the evidence and how |
| service_name | [string](#string) |  | Service_name is the name of service this evidence was collected from. For example, &#34;S3&#34; or &#34;GitLab&#34; |
| entity_type | [string](#string) |  | Entity_type specifies the row type and should correspond to a ServiceEntity. An entity_type typically represents a specific configurable entity such as AWS ECS &#34;Cluster&#34;. |
| sources | [Source](#receptor_v1-Source) | repeated | Sources are raw service provider API requests and responses used to generate this evidence. The raw API requests and responses serve as proof the evidence correlates to real service configurations. |
| doc | [Document](#receptor_v1-Document) |  | Document is an unstructured evidence. |
| struct | [Struct](#receptor_v1-Struct) |  | Struct is a structured evidence. |
| service_account_id | [string](#string) |  | Service_account_id is the member account of the organization. For example, an AWS account number in an organization that has more than one account. |






<a name="receptor_v1-Finding"></a>

### Finding
Finding is a set of evidence(s) collected from a service provider account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_type | [string](#string) |  | Unique receptor identifier. A receptor is expected to report findings from only one service provider type. A stable identifier that represent the type of receptor reporting this finding. The identifier is a simple URL encoded string that includes an organization name and the service provider name. For example: &#34;trustero_gitlab&#34;. |
| service_provider_account | [string](#string) |  | The receptor&#39;s evidence source. |
| entities | [ServiceEntity](#receptor_v1-ServiceEntity) | repeated | Entities is a list of service entity configurations in the service provider account. |
| evidences | [Evidence](#receptor_v1-Evidence) | repeated | One or more evidence collected by a typical receptor scan. |






<a name="receptor_v1-JobResult"></a>

### JobResult
JobResult reports the result of a receptor request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tracer_id | [string](#string) |  | Tracer_id is used to track the progress of the receptor request. |
| command | [string](#string) |  | Command is the receptor request that completed. One of &#34;verify&#34;, &#34;scan&#34;, or &#34;discover&#34; |
| result | [string](#string) |  | Result is receptor request result. One of &#34;success&#34;, &#34;fail&#34;, or &#34;error&#34;. |
| receptor_object_id | [string](#string) |  | Receptor_object_id is Trustero&#39;s receptor record identifier. |
| exceptions | [string](#string) |  | Exceptions contain information about the error like permission missing for the credentials provided. |






<a name="receptor_v1-ReceptorConfiguration"></a>

### ReceptorConfiguration
ReceptorConfiguration contains a configurations a receptor needs to access a service provider account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_object_id | [string](#string) |  | Receptor_object_id is Trustero&#39;s receptor record identifier. |
| credential | [string](#string) |  | Credential required to access a service provider for report finding and discover services purposes. |
| config | [string](#string) |  | Config holds additional receptor configuration to access a service provider account. |
| service_provider_account | [string](#string) |  | Service_provider_account is the service provider account name. |
| model_id | [string](#string) |  | Model_id is the receptor model id |






<a name="receptor_v1-ReceptorOID"></a>

### ReceptorOID
ReceptorOID is Trustero&#39;s receptor record identifier.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_object_id | [string](#string) |  | Receptor_object_id is the string representation of a Trustero persistent record. |






<a name="receptor_v1-ReportChunk"></a>

### ReportChunk



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_type | [string](#string) |  |  |
| service_provider_account | [string](#string) |  |  |
| data_chunk | [bytes](#bytes) |  | The chunk of the file/data (may contain either bytes or metadata) |
| file_name | [string](#string) |  | Metadata or additional optional information for the stream |
| mime_type | [string](#string) |  |  |
| is_final_chunk | [bool](#bool) |  | Indicates if this is the final chunk in the stream |






<a name="receptor_v1-ReportResponse"></a>

### ReportResponse
Server response after receiving the whole streamed report


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [string](#string) |  | Could use &#34;ok&#34; or other status information |
| file_path | [string](#string) |  | Optional: path where the server saved the chunked report |






<a name="receptor_v1-Row"></a>

### Row
Row is a row of structured data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity_instance_id | [string](#string) |  | Entity_instance_id of a discovered entity instance. For example, an AWS &#34;ECS&#34; cluster UUID or GitLab &#34;repository&#34; ID. |
| cols | [Row.ColsEntry](#receptor_v1-Row-ColsEntry) | repeated | Cols are columns of the row in column name to value pairs. All rows in a struct must have the same column names and corresponding value types. In addition, one of the key-value pair in the cols map must be the entity_instance_id, a unique instance of this row&#39;s Struct.entity_type. |






<a name="receptor_v1-Row-ColsEntry"></a>

### Row.ColsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Value](#receptor_v1-Value) |  |  |






<a name="receptor_v1-ServiceEntities"></a>

### ServiceEntities
ServiceEntities are configurable entities within a service provider account.  For example, service entities in
an AWS account include S3 buckets, ECS clusters, RDS database instances, etc.  The boundary of a service entity
depends on how a receptor represents it&#39;s findings.  Each service instance_id should be associated with at least
one Evidence.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| receptor_type | [string](#string) |  | Receptor_type is a unique receptor type. A stable string identifier that represent the type of receptor reporting this finding. The identifier is a simple URL encode string that includes the organization name and a service provider name. For example &#34;trustero_gitlab&#34;. |
| service_provider_account | [string](#string) |  | Service_provider_account is the service provider account where the services are configured in. |
| entities | [ServiceEntity](#receptor_v1-ServiceEntity) | repeated | Entities is a list of service instances configured in the service provider account. |






<a name="receptor_v1-ServiceEntity"></a>

### ServiceEntity
ServiceEntity is a discovered service entity instance which represents a configurable entity provided by the
service.  For example, an AWS ECS cluster or a GitLab repository.  These service entities are associated with
detailed configurations collected by Receptor.Reported() calls.  Service entities allows Trustero to filter
collected evidence that aren&#39;t relevant to specific audit contexts.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_name | [string](#string) |  | Service_name of the entity source. This is a Trustero assigned identifier for a known service such as &#34;GitLab&#34; or AWS &#34;ECS&#34;. |
| entity_type | [string](#string) |  | Entity_type is a service configurable object type such as a GitLab &#34;repository&#34; or AWS ECS &#34;cluster&#34;. The entity_instance_name and entity_instance_id must represent an instance of the subtype. For example, &#34;Java 1.5&#34; maybe a valid GitLab repository name or &#34;Elastic front end cluster&#34; maybe a valid AWS ECS cluster name. |
| entity_instance_name | [string](#string) |  | Entity_instance_name of a discovered service entity instance. For example, an AWS ECS cluster name or a GitLab repository name. Entity_instance_name of an entity may change for a given entity instance but it&#39;s entity_instance_id is stable. |
| entity_instance_id | [string](#string) |  | Entity_instance_id of a discovered entity instance. For example, an AWS ECS cluster UUID or GitLab repository ID. |
| service_account_id | [string](#string) |  | Service_account_id is the member account of the organization. For example, an AWS account number in a an organization that has more than one account. |






<a name="receptor_v1-Source"></a>

### Source
Source is the raw service provider API request and response.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| raw_api_request | [string](#string) |  | Raw_api_request is the raw API request used to generate this evidence. |
| raw_api_response | [string](#string) |  | Raw_api_response is the raw API response used to generate this evidence. |






<a name="receptor_v1-Struct"></a>

### Struct
Struct is a structured evidence defined in tabular form.  Each struct typically represent a service type (see
Evidence message definition).  Each struct consists of rows of data.  Each row typically represent a service
instance and its configurations.  A row  contains column name and column value pairs.  All rows in a struct
must have the same column name-value pairs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rows | [Row](#receptor_v1-Row) | repeated | Rows of key-value pairs. Each row typically represents the configuration of a service instance or an data type such as a member of GitLab group. |
| col_display_names | [Struct.ColDisplayNamesEntry](#receptor_v1-Struct-ColDisplayNamesEntry) | repeated | Col_display_names is a map of row column name to it&#39;s corresponding display name. Display names are used by the user interface to render a field in a struct&#39;s rows. |
| col_display_order | [string](#string) | repeated | Col_display_order is an ordered list of row column names. The order of the column names are used by the user interface to render the column order of a struct&#39;s rows. |
| col_tags | [Struct.ColTagsEntry](#receptor_v1-Struct-ColTagsEntry) | repeated | Col_tags is a map of column tag to row column name. These tags are used by Trustero&#39;s evidence testing interface to identify relevant evidence. |






<a name="receptor_v1-Struct-ColDisplayNamesEntry"></a>

### Struct.ColDisplayNamesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="receptor_v1-Struct-ColTagsEntry"></a>

### Struct.ColTagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="receptor_v1-Value"></a>

### Value
Value is a [Struct.row.col] column value.  Value types can be simple protobuf scalar or [google.proto.Timestamp].


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
Receptor service, or a Trustero client application, collects findings supporting the use of services from a
service provider account.  For example, AWS is a service provider, AWS account is a service provider account,
and S3 is a service.  Trustero associates collected evidence to business controls in support of the business&#39;
stated procedures.  Evidences are organized into a finding.  Each evidence is associated with a service entity
and contains the service&#39;s configuration information. For example, an AWS S3 bucket and its configuration.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Verified | [Credential](#receptor_v1-Credential) | [.google.protobuf.Empty](#google-protobuf-Empty) | Verified reports whether the provided credential is a valid service provider credential for purpose of discovering service entities and reporting findings. This rpc call is typically made as callback by a receptor to Trustero from a check-credential receptor request. |
| GetConfiguration | [ReceptorOID](#receptor_v1-ReceptorOID) | [ReceptorConfiguration](#receptor_v1-ReceptorConfiguration) | GetConfiguration for the receptor and service provider credential using the provided receptor object identifier. This rpc call is typically made as a callback by a receptor prior to making a report findings or discover service entities receptor request. |
| Discovered | [ServiceEntities](#receptor_v1-ServiceEntities) | [.google.protobuf.StringValue](#google-protobuf-StringValue) | Discovered reports known service entities. A receptor or a Trustero client application reports its known service entities on request. A service entity is a configurable asset of a service such as an S3 bucket. This call returns a string value discovery ID or an error. |
| Report | [Finding](#receptor_v1-Finding) | [.google.protobuf.StringValue](#google-protobuf-StringValue) | Report a finding to Trustero. A receptor or a Trustero client application reports its findings to Trustero on a periodic basis. This call returns a string value collection ID or an error. |
| Notify | [JobResult](#receptor_v1-JobResult) | [.google.protobuf.Empty](#google-protobuf-Empty) | Notify Trustero a long running report finding or discover service entities receptor-request has completed. JobResult contains information about the receptor-request and it&#39;s corresponding result. |
| SetConfiguration | [ReceptorConfiguration](#receptor_v1-ReceptorConfiguration) | [.google.protobuf.Empty](#google-protobuf-Empty) | SetConfiguration reports the configuration for receptors that need extra configuration to access a service. This call is typically made as a callback by a receptor after credential verification. |
| StreamReport | [ReportChunk](#receptor_v1-ReportChunk) stream | [ReportResponse](#receptor_v1-ReportResponse) |  |

 



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

