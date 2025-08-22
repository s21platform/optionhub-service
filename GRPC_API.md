# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/optionhub.proto](#api_optionhub-proto)
    - [AttributeMetadata](#-AttributeMetadata)
    - [GetAttributesMetadataIn](#-GetAttributesMetadataIn)
    - [GetAttributesMetadataOut](#-GetAttributesMetadataOut)
    - [GetAttributesMetadataOut.AttributesMetadataEntry](#-GetAttributesMetadataOut-AttributesMetadataEntry)
  
    - [AttributeType](#-AttributeType)
    - [EntityType](#-EntityType)
  
    - [OptionhubService](#-OptionhubService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_optionhub-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/optionhub.proto



<a name="-AttributeMetadata"></a>

### AttributeMetadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attribute_id | [int64](#int64) |  |  |
| type | [AttributeType](#AttributeType) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| allowed_operators | [bytes](#bytes) |  |  |
| entity_attribute_id | [int64](#int64) |  |  |
| entity_type | [EntityType](#EntityType) |  |  |
| label | [string](#string) |  |  |
| is_required | [bool](#bool) |  |  |
| order_index | [int64](#int64) |  |  |
| visibility_rules | [bytes](#bytes) |  |  |






<a name="-GetAttributesMetadataIn"></a>

### GetAttributesMetadataIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity_attribute_ids | [int64](#int64) | repeated |  |






<a name="-GetAttributesMetadataOut"></a>

### GetAttributesMetadataOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attributes_metadata | [GetAttributesMetadataOut.AttributesMetadataEntry](#GetAttributesMetadataOut-AttributesMetadataEntry) | repeated |  |






<a name="-GetAttributesMetadataOut-AttributesMetadataEntry"></a>

### GetAttributesMetadataOut.AttributesMetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [int64](#int64) |  |  |
| value | [AttributeMetadata](#AttributeMetadata) |  |  |





 


<a name="-AttributeType"></a>

### AttributeType


| Name | Number | Description |
| ---- | ------ | ----------- |
| STRING | 0 |  |
| NUMBER | 1 |  |
| DATE | 2 |  |
| OPTION | 3 |  |
| MULTISELECT | 4 |  |



<a name="-EntityType"></a>

### EntityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER | 0 |  |
| SOCIETY | 1 |  |
| EVENT | 2 |  |
| COMMUNITY | 3 |  |
| MATERIAL | 4 |  |


 

 


<a name="-OptionhubService"></a>

### OptionhubService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAttributesMetadata | [.GetAttributesMetadataIn](#GetAttributesMetadataIn) | [.GetAttributesMetadataOut](#GetAttributesMetadataOut) |  |

 



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

