SFS
===================
In-browser simple Linux-like file system.

User can interact with the system through a chat-like interface.

![](/assets/fs.png)


System operations
-------------

Note: all paths used to navigate must be from root ('/')

- List files/folders: 
```
ls /path/to/folder
```
- Show file data: 
```
cat /path/to/file
```
- Create new file/folder: 
```
cr /path/to/file [optional_data]
```
- Find files/folders: 
```
find file_name_substring
```
- Remove file/folder: 
```
rm /path/to/delete/1 [/path/to/delete/2, /path/to/delete/3]
```
- Move file/folder: 
```
mv /path/to/file/or/folder new_name [optional_data]
```


TODO:
-------------
- [ ] Navigate around the file system
- [ ] Password authentication
- [ ] Multi-user support

### Prerequiste
- go 1.17

### Run localy
```
cd go-backend
make local-db
make setup
make run
```