
hey:
	hey -h2 -cpus 2 https://127.0.0.1:8888/user?X-Access-Token=MTQ5ODQ5NzQ3N3x6NnItY0hDQmMyaXZrS2c0eTduRDNCNzhVbWhFQmVfR2pnOHVKWVhER2l5eXBJczFicG81V1F6RFA4Yz18F9YI7t_OGXU395L68OtlZ7eUizs3x9bwVYtbQvyYcCY=

flatcAPI:
	flatc -g api/structsz/schema_MiddCache.fbs 

flatcJS:
	cd public/js/ && flatc --js --no-js-exports middCache.fbs 


flatcTest:
	flatc -g api/structsz/schema_objs.fbs