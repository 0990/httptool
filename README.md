# httptool
web tool,such as create qrcode,file server,etc

## install
Docker compose
```yaml
version: '3'
 
services:
  httptool:
    image: 0990/httptool:latest
    environment:
      - Advertised_Addr=http://127.0.0.1
    ports:
      - "80:80"
    volumes:
      - xxx:/httptool/file_dir
      - xxx:/httptool/qrcode_dir
```
## web html tool
```http request
/tool.html
```

## qrcode 
create qrcode
```http request
/qrcode/create POST
Content-Type: application/json
```
Body:
```json
{
  "content":"https://www.google.com"
}
```
Response:
```json
{"code":0,"message":"","picUrl":"http://127.0.0.1/qrcode/qrcode20230906143105.png"}
```
picUrl is the qrcode url,"http://127.0.0.1" can configure by env **Advertised_Addr**
generated png file dir can configure by volumes **/httptool/qrcode_dir**

## file server
```http request
/file
```
generated png file dir can configure by volumes **/httptool/file_dir**



