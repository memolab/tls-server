

flatcAPI:
	flatc -g api/structsz/schema_MiddCache.fbs 

flatcJS:
	cd public/js/ && flatc --js --no-js-exports middCache.fbs 