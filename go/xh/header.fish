set XH_URL   (cat url.var)
set XH_TOKEN (cat token.var)
set XH_AUTH  -A bearer -a $XH_TOKEN

function request -a method url
  set -l xh_additional_args $argv[3..-1]
  echo $method $XH_URL$url

  xh $method $XH_URL$url $xh_additional_args
end

# this is not really here for functionality
# rather to make the files more understandable
function request_raw -a method url
  set -l xh_additional_args $argv[3..-1]
  set -l req (xh $method $XH_URL$url $xh_additional_args)
  echo $req
end

function check -a string
  set -l has_error (echo $string | jq -e 'has("error") or has("error_message") or .token == null' >/dev/null 2>&1; echo $status)
end

function set_var_file -a string name 
  set -l ret $res | jq -r ".$name" 
  set has_error (check $res)
  if [ "$ret" != "null" ] 
    if [ "$has_error" = "0" ]
      echo "$ret" > "$name.var" 
      return
    end
    echo error! 
  end
end
