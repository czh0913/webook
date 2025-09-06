local key = KEYS[1]
-- 用户输入的 code
local expectedCode = ARGV[1]
local cntKey = key..":cnt"
local code = redis.call("get",key)
local cnt = tonumber(redis.call("get",  cntKey))


if cnt <= 0 then
-- 说明用户一直出错
    return -1
elseif expectedCode == code then
    redis.call("set",cntKey,-1)
    return 0
else
    redis.call("decr",cntKey)
    return -2
end
