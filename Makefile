torch:
	go-torch -t=5 -u="http://127.0.0.1:8888" -p > profile.svg

pprof:
	go tool pprof --seconds=5 http://127.0.0.1:8888/debug/pprof/profile


hey:
	hey -h2 -cpus 2 -H "Content-Type: text/plain" -H "Connection: keep-alive" http://127.0.0.1:8888/
heyUsr:
	hey -h2 -cpus 2 -H "Content-Type: application/json" -H "Connection: keep-alive" http://127.0.0.1:8888/user?X-Access-Token=MTQ5ODQ5NzQ3N3x6NnItY0hDQmMyaXZrS2c0eTduRDNCNzhVbWhFQmVfR2pnOHVKWVhER2l5eXBJczFicG81V1F6RFA4Yz18F9YI7t_OGXU395L68OtlZ7eUizs3x9bwVYtbQvyYcCY=
heyUsrOne:
	hey -h2 -cpus 2 -c 1 -n 1 -H "Content-Type: application/json" -H "Connection: keep-alive" http://127.0.0.1:8888/user?X-Access-Token=MTQ5ODQ5NzQ3N3x6NnItY0hDQmMyaXZrS2c0eTduRDNCNzhVbWhFQmVfR2pnOHVKWVhER2l5eXBJczFicG81V1F6RFA4Yz18F9YI7t_OGXU395L68OtlZ7eUizs3x9bwVYtbQvyYcCY=
heyUsr2:
	hey -h2 -cpus 2 -H "Content-Type: text/plain" -H "Connection: keep-alive" http://127.0.0.1:8888/user2?X-Access-Token=MTQ5ODQ5NzQ3N3x6NnItY0hDQmMyaXZrS2c0eTduRDNCNzhVbWhFQmVfR2pnOHVKWVhER2l5eXBJczFicG81V1F6RFA4Yz18F9YI7t_OGXU395L68OtlZ7eUizs3x9bwVYtbQvyYcCY=

heySignin:
	hey -h2 -cpus 2 -H "Content-Type: application/json" -H "Connection: keep-alive" -H "X-Forwarded-For: 196.221.105.74" https://127.0.0.1:8888/signin


flatcAPI:
	flatc -g api/structsz/schema_MiddCache.fbs 

flatcJS:
	cd public/js/ && flatc --js --no-js-exports middCache.fbs 