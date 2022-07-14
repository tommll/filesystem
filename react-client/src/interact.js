// require('dotenv').config();

export const sendToServer =  async (args) => {
  // var parts = JSON.parse(args)['command'].split(' ')
  var parts = args.split(' ')
  var cmd = parts[0]

  var ans = "";
  var url = 'https://nameless-taiga-69302.herokuapp.com/fs/' + cmd
  var data;

  switch (cmd){
    case 'ls':
      if (parts.length == 2){
        data = {"path": parts[1]}
      } else {
        data = {"path": ""}
      }
      break
    case 'cr':
        if (parts.length == 2){
          data = {"path": parts[1], "data": ""}
        } else {
          data = {"path": parts[1], "data": parts.slice(2).join(' ')}
        }
        break
    case 'cat':
        data = {"path": parts[1]}
        break
    case 'find':
        if (parts.length == 2){
          data = {"name": parts[1], "path": ""}
        } else {
          data = {"name": parts[1], "path": parts[2]}
        }
        break
    case 'up':
        if (parts.length == 3){
          data = {"path": parts[1], "name": parts[2], "data": ""}
        } else {
          data = {"path": parts[1], "name": parts[2], "data": parts.slice(3).join(' ')}
        }
        break
    case 'mv':
        data = {"src_path": parts[1], "dest_path": parts[2]}
        break
    case 'rm':
        data = {"paths": parts.slice(1)}
        break
    default:
      console.log('invalid command: ', cmd)
      return ""
  }

  console.log(`sending url: ${url}, data: ${JSON.stringify(data)}`)
  
  var requestOptions = {
    method: 'POST',
    body: JSON.stringify(data)
  }

  try{
    var res = await fetch(url, requestOptions)
    .then(response => response.json())
    .then(result => {
      var stt = result.return_code
      console.log(`response type: ${typeof result}: ${result}`)

      if (stt == 1){
        if (result.data != null){
          var data = result.data
          if (data.items != null){
            ans = parseList(Object.values(data.items))
          } else {
            ans = data.toString()
          }
        } else {
          ans = ""
        }
        console.log(`message: ${result.message.toString()}}`)

      } else {
        console.log(`err: ${result.message}`)
        ans = `Error: ${result.message.toString()}`
      }
      
  })
  } catch (error) {
    console.log('error:', error)
    return "Error: something wrong with the network"
  }
  
  return ans;
}

export function checkInput(data) {
  var parts = data.split(' ')

  if (parts[0] === "ls" ){
    return true
  } else if (parts[0] === "cat" ){
    if (parts.length != 2){
      return false
    }
    return true
  } else if (parts[0] === "cr" ){
    if (parts.length < 2){
      return false
    }
    return true
  } else if (parts[0] === "find"){
    if (parts.length < 2){
      return false
    }
    return true
  } else if (parts[0] === "up"){
    if (parts.length < 3){
      return false
    }
    return true
  } else if (parts[0] === "mv"){
    if (parts.length != 3){
      return false
    }
    return true
  } else if (parts[0] === "rm"){
    if (parts.length < 2){
      return false
    }
    return true
  } 
  return false
}

// Input: a list of files/folders
// A string display files and folders under a directory like in terminal
// 
// ~/tung$
// file1   folder2   folder4
var LINE_LENGTH = 50;

function parseList(list) {
  var res = "";
  var currentLength = 0;

  for (var i = 0; i < list.length; i++){
    if (list[i].length + currentLength + 6 > LINE_LENGTH){
      currentLength = list[i].length
      res += "\n" + list[i]
    } else {
      currentLength += list[i].length + 6
      if (i != 0){
        res += "      " + list[i]
      } else {
        res += list[i]
      }
    }
  }

  return res
}