-- Natsuki HTTP Driver
-- Version 1.0 - July 5th, 2022
local NatsukiURL = "<URL Where Natsuki is self-hosted>"
local NatsukiJWT = "<Your Self Signed JWT Token>"

local HttpService = game:GetService("HttpService")
local Natsuki = {
	["RedisURL"] = string.format("%s/redis", NatsukiURL),
	["PostgresURL"] = string.format("%s/sql", NatsukiURL),
	["Headers"] = {
		["Content-Type"] = "application/json",
		["Authorization"] = NatsukiJWT,
		["RBX-Server-ID"] = game.JobId,
	},
}

type HTTPResponse = {
	Success: boolean,
	StatusCode: number,
	StatusMessage: string,
	Headers: { [string]: string },
	Body: string,
}

-- Create a Redis Command and run it
-- Example:
-- Command = {"HSET", "FOO", "BAR", "BAZ"}
-- Result = 1
function Natsuki:RedisAsync(Command: { any }): (boolean, any)
	-- Convert Commands to Strings
	for i, Cmd in Command do
		Command[i] = tostring(Cmd)
	end

	-- Send Request to Natsuki API
	local Ran: boolean, Response: HTTPResponse = pcall(function()
		HttpService:RequestAsync({
			["Method"] = "POST",
			["Url"] = Natsuki.RedisURL,
			["Headers"] = Natsuki.Headers,
			["Body"] = HttpService:JSONEncode(Command),
		})
	end)

	-- Encountered HTTP Error?
	if Ran == false then
		return false, string.format("HTTP Error: %s", Response)
	end

	-- Parse JSON Body
	local Ran: boolean, Body: any = pcall(function()
		return HttpService:JSONEncode(Response.Body).result
	end)

	-- Encountered JSON Parsing Error?
	if Ran == false then
		return false, string.format("JSON Parse Error: %s", Body)
	end

	-- Encountered Server Error?
	if Body.error then
		return false, string.format("Error Code: %s, Error Message: %s", Body.code, Body.error)
	end

	-- Return Result
	return true, Body
end

-- Create a Postgres Query and run it (Natsuki will automatically escape it to prevent SQL Injection)
-- Example:
-- Query = "INSERT INTO monke.history(id, timestamp, player_id, stage_id) VALUES (DEFAULT, DEFAULT, $1, $2) RETURNING id;"
-- Arguments = 41537575, 8008135
-- Result = true, {10}
function Natsuki:PostgresAsync(Query: string, Arguments: { any }): (boolean, any)
	-- Send Request to Natsuki API
	local Ran: boolean, Response: HTTPResponse = pcall(function()
		HttpService:RequestAsync({
			["Method"] = "POST",
			["Url"] = Natsuki.PostgresURL,
			["Headers"] = Natsuki.Headers,
			["Body"] = HttpService:JSONEncode({
				["query"] = Query,
				["args"] = Arguments,
			}),
		})
	end)

	-- Encountered HTTP Error?
	if Ran == false then
		return false, string.format("HTTP Error: %s", Response)
	end

	-- Parse JSON Body
	local Ran: boolean, Body: any = pcall(function()
		return HttpService:JSONEncode(Response.Body).result
	end)

	-- Encountered JSON Parsing Error?
	if Ran == false then
		return false, string.format("JSON Parse Error: %s", Body)
	end

	-- Encountered Server Error?
	if Body.error then
		return false, string.format("Error Code: %s, Error Message: %s", Body.code, Body.error)
	end

	-- Return Result
	return true, Body
end

-- hello its me natsuki
return Natsuki
