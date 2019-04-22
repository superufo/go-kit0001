main.exe  -consul.host 47.112.111.171 -consul.port 8500  -service.host 110.235.246.150   -service.port 9550          
main.exe   -consul.host 47.112.111.171  -consul.port 8500


zipkin.exe -consul.host 47.112.111.171 -consul.port 8500  --zipkin.url http://47.112.111.171:9411/zipkin/api/v2/spans    

register.exe -consul.host 47.112.111.171 -consul.port 8500  -service.host 110.235.246.150   -service.port 9550  
--zipkin.url http://47.112.111.171:9411/zipkin/api/v2/spans    


http://47.112.111.171:8181/hystrix
http://127.0.0.1:8880/arithmetic/calculate/Add/9851/1245821