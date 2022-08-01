# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [receptor_v1/receptor.proto](#receptor_v1_receptor-proto)
    - [Credential](#receptor_v1-Credential)
    - [Credential.CredentialEntry](#receptor_v1-Credential-CredentialEntry)
    - [JobResult](#receptor_v1-JobResult)
    - [ReceptorConfiguration](#receptor_v1-ReceptorConfiguration)
    - [ReceptorOID](#receptor_v1-ReceptorOID)
    - [Services](#receptor_v1-Services)
    - [Services.Service](#receptor_v1-Services-Service)
  
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
| credential | [Credential.CredentialEntry](#receptor_v1-Credential-CredentialEntry) | repeated | The service provider credential being verified. |
| is_credential_valid | [bool](#bool) |  | Report whether the service provider credential provided in this message is valid for report findings or discover services request. |
| message | [string](#string) |  | Reason for why the service provider credential in this message is invalid. |






<a name="receptor_v1-Credential-CredentialEntry"></a>

### Credential.CredentialEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






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
| receptor_type | [string](#string) |  | Unique receptor type. A stable string identifier that represent the type of receptor reporting this finding. The identifier is akin to a fully qualified Go package name or a Java class name. For example, &#34;github.com/trustero/receptor/gitlab&#34;. @required |
| service_provider_account | [string](#string) |  | The service provider of this list of services. @required |
| services | [Services.Service](#receptor_v1-Services-Service) | repeated | A list of service instances. @required |






<a name="receptor_v1-Services-Service"></a>

### Services.Service
A service instance definition.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the service. For example, &#34;ECS&#34;. @required |
| instance_id | [string](#string) |  | Unique service ID. For example, ECS&#39;s UUID. @required |





 

 

 


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

