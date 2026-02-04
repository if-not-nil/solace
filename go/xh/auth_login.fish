source header.fish

set res (request_jq post "/auth/login" name=user password=password $auth)

set -l token $res | jq -r '.token' 
echo $res | jq
if [ "$token" != "null" ] 
  echo "$token" > token.var 
end
