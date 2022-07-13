--[[
	Name: Natsuki Driver
	Developer: @bakonpancakz
	Release Date: July 12th 2022
	Version: 2.0
	---
	Function Responses:
		All Database Command/Query Functions currently return two values:
		
		[#1] Boolean:
		Whether or not the function completed succesfully

		[#2] string OR any:
		If Value #1 is true, then whatever your function returns, it has returned.
	---
	Script Attributes:
		Imagine these as the "environment variables" for this client, set these as attributes the ModuleScript you created.
		We recommend you place this ModuleScript in either ServerStorage or ServerScriptService whichever you prefer so you can easily require() it.
		| Name       | Type    |       Default       | Description                                             |
		| :--------- | :------ | :-----------------: | :------------------------------------------------------ |
		| EnableLogs | boolean |        true         | Whether or not to log to console                        |
		| Token      | string  |  <Your JWT Token>   | Your Authorization Token, self signed using your secret |
		| ProUrl     | string  |   <Your Host URL>   | Your Production URL for Natsuki                         |
		| DevUrl     | string  | http://localhost:80 | Your Development URL for Natsuki                        |
	---
]]
local HttpService = game:GetService("HttpService")
local Natsuki = {
	["JobId"] = game.JobId,
	["EnableLogs"] = script:GetAttribute("EnableLogs") or true,
	["URL"] = script:GetAttribute("ProUrl") or "<Your Host URL>",
	["JWT"] = script:GetAttribute("Token") or "<Your JWT Token>",
}

-- Run Special Options when in Roblox Studio
if game:GetService("RunService"):IsStudio() then
	-- Set this to your localhost address and port
	Natsuki.URL = script:GetAttribute("DevUrl") or "http://localhost:80"

	-- Set JobId as it will be "" in Studio which will a 401: Unauthorized Error
	Natsuki.JobId = HttpService:GenerateGUID(false)

	-- Notify we are running in Development Mode! So start your Local Server
	warn(string.format("[NATSUKI] Running in Development Mode!"))
end

--[[
	Create a Redis Command and execute it
	If you are executing a redis command then expect it to return what running it in the cli would return
	You must parse it yourself though, because im not that smart.

	Examples:
	SET foo bar => string: OK
	GET foo bar => string: bar
	SADD foo bar baz => string: 2
	SMEMBERS => table: {string: bar, string: baz}
]]
function Natsuki:RedisAsync(Command: { any }): (boolean, string)
	-- Convert Commands to Strings
	for Index, CommandValue in Command do
		Command[Index] = tostring(CommandValue)
	end

	local Ran, Result: any = pcall(function()
		-- [1] Send HTTP Request
		local HTTPRequest = HttpService:RequestAsync({
			["Method"] = "POST",
			["Body"] = HttpService:JSONEncode(Command),
			["Url"] = string.format("%s/redis", Natsuki.URL),
			["Headers"] = {
				["Content-Type"] = "application/json",
				["RBX-Server-Id"] = Natsuki.JobId,
				["Authorization"] = Natsuki.JWT,
			},
		})

		-- [2] Parse JSON Body
		local Body = HttpService:JSONDecode(HTTPRequest.Body)

		-- [3] Return Client Error (if any)
		if Body.error then
			return string.format("%s (%c)", Body.message, Body.error)
		end

		-- [4] Return Result
		return Body.result
	end)

	-- Log Request to Console
	-- [NATSUKI:RDB] SADD tracker:players:unique 41537575 => 1
	-- [NATSUKI:RDB] SADD tracker:players:unique 41537575 => example redis error (1273)
	if Natsuki.EnableLogs then
		print(
			string.format(
				"[NATSUKI:RDB] %s => %s",
				table.concat(Command, " "),
				string.sub(HttpService:JSONEncode(Result or "N/A"), 0, 64)
			)
		)
	end

	return Ran, Result
end

--[[
	Create and Execute a Postgres Query via Natsuki:
	Example:
	Ran, Results = Natsuki:PostgresAsync(
		Query = "INSERT INTO monke.history(player_id, stage_id) VALUES ($1, $2), ($3, $4) RETURNING id;",
		Arguments = {41537575, 8008135, 12345678, 8008135},
		Schema = {"id"}
	)
	Ran = true
	Results = {
		{id: 6},
		{id: 7}
	}
]]
function Natsuki:PostgresAsync(Query: string, Arguments: { any }, Schema: { string }): (boolean, { { [string]: any } })
	local Ran, Result: string | { { [string]: any } } = pcall(function()
		-- [1] Send HTTP Request
		local HTTPRequest = HttpService:RequestAsync({
			["Method"] = "POST",
			["Url"] = string.format("%s/sql", Natsuki.URL),
			["Body"] = HttpService:JSONEncode({
				["query"] = Query,
				["args"] = Arguments,
			}),
			["Headers"] = {
				["Content-Type"] = "application/json",
				["RBX-Server-Id"] = Natsuki.JobId,
				["Authorization"] = Natsuki.JWT,
			},
		})

		-- [2] Parse JSON Body
		local Body = HttpService:JSONDecode(HTTPRequest.Body)

		-- [3] Return Client Error (if any)
		if Body.error then
			return string.format("%s (%c)", Body.message, Body.error)
		end

		-- [4] Create Table with Schema
		local Results = {}
		for _, Row: { any } in Body.result do
			local RowTable = {}
			for index: number, Value: any in Row do
				RowTable[Schema[index]] = Value
			end
			table.insert(Results, RowTable)
		end

		-- [5] Return Row(s)
		return Results
	end)

	-- Log Request to Console
	-- [NATSUKI:SQL] SELECT * FROM monke.history; => [{id: 1, timestamp = "2022-07-04T20:11:45.871965Z", player_id = 41537575, stage_id = 2}]
	-- [NATSUKI:SQL] INSERT INTO monke.history(player_id, stage_id) VALUES ($1, $2) RETURNING id; => [{id: 13}]
	-- [NATSUKI:SQL] INSERT INTO monke.history(player_id, stage_id) VALUES ($1, $2); => example postgres error (1273)
	if Natsuki.EnableLogs then
		print(
			string.format(
				"[NATSUKI:SQL] %s (%s) => %s",
				Query,
				table.concat(Arguments, " "),
				string.sub(HttpService:JSONEncode(Result or "N/A"), 0, 64)
			)
		)
	end

	return Ran, Result
end

-- hello i am natsuki!!! 🧁
return Natsuki
