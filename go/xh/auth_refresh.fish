source header.fish

set res (request_raw get "/auth/refresh" name=user password=password $auth)

echo $res | jq
set_var_file $res token
