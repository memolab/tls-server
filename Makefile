
hey:
	hey -h2 -cpus 2 https://127.0.0.1:8888/user?X-Access-Token=MTUwMDYzMjE0NnxFQnZrc0xFMTFlUVlaWGVkWW1RSWxPQWRYRjA1bktveU84SXJsc0hLcmwzdXZjTllfUkJ0TXc9PXwnKxYKx9Tew8S2qKhZ_VP07LsPQvvyiM9C2tfvRw_rPg==


flatcAPI:
	flatc -g api/structsz/schema_MiddCache.fbs

flatcJS:
	cd public/js/ && flatc --js --no-js-exports middCache.fbs
