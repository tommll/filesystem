SFS
===================
In-browser simple Linux-like file system

![](/assets/fs.png)


System operations
-------------
- List all file and folders in a directory 
- Show data of a file 
- Create a new file/folder
- Find all files and folders that have name matches a substring
- Remove a file/folder 
- Move a file/folder to another directory

- Note: all paths used to navigate must be from root ('/')

Minimal interface
-------------
User can interact with the system through a chat-like interface 
- List files/folders: ls /path/to/folder
- Show file data: cat /path/to/file
- Create new file/folder: cr /path/to/file [optional_data]
- Find files/folders: find file_name_substring
- Remove file/folder : rm /path/to/delete/1 [/path/to/delete/2, /path/to/delete/3]
- Move file/folder : mv /path/to/file/or/folder new_name [optional_data]



Non-functional features:
-------------
- Security: password-based authentication


### Run localy
```
cd go-backend
make local-db
make setup
make run
```