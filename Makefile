gen_auth_service: proto/auth/auth_service.proto
	protoc --proto_path=./proto/auth  --go_opt=paths=source_relative --go_out=./internal/gen/proto/auth --go-grpc_out=./internal/gen/proto/auth  --go-grpc_opt=paths=source_relative  auth_service.proto
gen_user_manegement: proto/user_manegement/user_manegement.proto 
	protoc --proto_path=./proto/user_manegement  --go_opt=paths=source_relative --go_out=./internal/gen/proto/user_manegement --go-grpc_out=./internal/gen/proto/user_manegement  --go-grpc_opt=paths=source_relative  user_manegement.proto 