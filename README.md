# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [optionhub_service.proto](#optionhub_service-proto)
    - [AddIn](#-AddIn)
    - [AddOut](#-AddOut)
    - [DeleteByIdIn](#-DeleteByIdIn)
    - [DeleteByIdOut](#-DeleteByIdOut)
    - [GetAllIn](#-GetAllIn)
    - [GetAllOut](#-GetAllOut)
    - [GetByIdIn](#-GetByIdIn)
    - [GetByIdOut](#-GetByIdOut)
    - [SetByIdIn](#-SetByIdIn)
    - [SetByIdOut](#-SetByIdOut)
  
    - [OptionhubService](#-OptionhubService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="optionhub_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## optionhub_service.proto



<a name="-AddIn"></a>

### AddIn
message request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  | value that will be added |






<a name="-AddOut"></a>

### AddOut
message response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  | id of added value |
| value | [string](#string) |  |  |






<a name="-DeleteByIdIn"></a>

### DeleteByIdIn
message request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  | value with such id will be marked as inactive |






<a name="-DeleteByIdOut"></a>

### DeleteByIdOut
message response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | flag that says if deletion was successful |






<a name="-GetAllIn"></a>

### GetAllIn







<a name="-GetAllOut"></a>

### GetAllOut
message response for all active values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [GetByIdOut](#GetByIdOut) | repeated |  |






<a name="-GetByIdIn"></a>

### GetByIdIn
message request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  | user&#39;s id |






<a name="-GetByIdOut"></a>

### GetByIdOut
message response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| value | [string](#string) |  |  |






<a name="-SetByIdIn"></a>

### SetByIdIn
message request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |






<a name="-SetByIdOut"></a>

### SetByIdOut
message response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  | id of value that will be changed |
| value | [string](#string) |  |  |





 

 

 


<a name="-OptionhubService"></a>

### OptionhubService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetOsById | [.GetByIdIn](#GetByIdIn) | [.GetByIdOut](#GetByIdOut) |  |
| GetAllOs | [.GetAllIn](#GetAllIn) | [.GetAllOut](#GetAllOut) |  |
| AddOs | [.AddIn](#AddIn) | [.AddOut](#AddOut) |  |
| SetOsById | [.SetByIdIn](#SetByIdIn) | [.SetByIdOut](#SetByIdOut) |  |
| DeleteOsById | [.DeleteByIdIn](#DeleteByIdIn) | [.DeleteByIdOut](#DeleteByIdOut) |  |
| GetWorkPlaceById | [.GetByIdIn](#GetByIdIn) | [.GetByIdOut](#GetByIdOut) |  |
| GetAllWorkPlace | [.GetAllIn](#GetAllIn) | [.GetAllOut](#GetAllOut) |  |
| AddWorkPlace | [.AddIn](#AddIn) | [.AddOut](#AddOut) |  |
| SetWorkPlaceById | [.SetByIdIn](#SetByIdIn) | [.SetByIdOut](#SetByIdOut) |  |
| DeleteWorkPlaceById | [.DeleteByIdIn](#DeleteByIdIn) | [.DeleteByIdOut](#DeleteByIdOut) |  |
| GetStudyPlaceById | [.GetByIdIn](#GetByIdIn) | [.GetByIdOut](#GetByIdOut) |  |
| GetAllStudyPlace | [.GetAllIn](#GetAllIn) | [.GetAllOut](#GetAllOut) |  |
| AddStudyPlace | [.AddIn](#AddIn) | [.AddOut](#AddOut) |  |
| SetStudyPlaceById | [.SetByIdIn](#SetByIdIn) | [.SetByIdOut](#SetByIdOut) |  |
| DeleteStudyPlaceById | [.DeleteByIdIn](#DeleteByIdIn) | [.DeleteByIdOut](#DeleteByIdOut) |  |
| GetHobbyById | [.GetByIdIn](#GetByIdIn) | [.GetByIdOut](#GetByIdOut) |  |
| GetHobbyPlace | [.GetAllIn](#GetAllIn) | [.GetAllOut](#GetAllOut) |  |
| AddHobby | [.AddIn](#AddIn) | [.AddOut](#AddOut) |  |
| SetHobbyById | [.SetByIdIn](#SetByIdIn) | [.SetByIdOut](#SetByIdOut) |  |
| DeleteHobbyById | [.DeleteByIdIn](#DeleteByIdIn) | [.DeleteByIdOut](#DeleteByIdOut) |  |
| GetSkillById | [.GetByIdIn](#GetByIdIn) | [.GetByIdOut](#GetByIdOut) |  |
| GetSkillPlace | [.GetAllIn](#GetAllIn) | [.GetAllOut](#GetAllOut) |  |
| AddSkill | [.AddIn](#AddIn) | [.AddOut](#AddOut) |  |
| SetSkillById | [.SetByIdIn](#SetByIdIn) | [.SetByIdOut](#SetByIdOut) |  |
| DeleteSkillById | [.DeleteByIdIn](#DeleteByIdIn) | [.DeleteByIdOut](#DeleteByIdOut) |  |
| GetCityById | [.GetByIdIn](#GetByIdIn) | [.GetByIdOut](#GetByIdOut) |  |
| GetCityPlace | [.GetAllIn](#GetAllIn) | [.GetAllOut](#GetAllOut) |  |
| AddCity | [.AddIn](#AddIn) | [.AddOut](#AddOut) |  |
| SetCityById | [.SetByIdIn](#SetByIdIn) | [.SetByIdOut](#SetByIdOut) |  |
| DeleteCityById | [.DeleteByIdIn](#DeleteByIdIn) | [.DeleteByIdOut](#DeleteByIdOut) |  |

 



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

