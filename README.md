## SFS
A simple version of Linux file system


### Core features:
#### *1. File system operations*
- List all file and folders in a directory 
- Show data of a file 
- Create a new file/folder
- Find all files and folders that have name matches a substring
- Remove a file/folder 
- Move a file/folder to another directory

- Note: all paths used to navigate must be from root ('/', './')

#### *2. Minimal interface*
- Currently, user can only interact with the system through http protocol
- List files/folders: ${host}/fs/ls
- Show file data: ${host}/fs/cat
- Create new file/folder: ${host}/fs/cr
- Find files/folders: ${host}/fs/find
- Remove file/folder : ${host}/fs/rm
- Move file/folder : ${host}/fs/mv

### Non-functional features:
- Security: password-based authentication