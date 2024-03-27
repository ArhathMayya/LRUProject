import logo from './logo.svg';
import './App.css';
import * as React from 'react';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import { useState } from 'react';
import axios from 'axios';
import Paper from '@mui/material/Paper';
import Alert from '@mui/material/Alert';

function App() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("")
  const [getValue, setGetValue] = useState("")
  const [cacheAlert, setCacheAlert] = useState("")
  const [getValueAlert, SetGetValueAlert] = useState("")
  const [error, SetError] = useState(false)

  const handleSendButton = () =>{
    
    console.log(key, value)

    axios.post('http://localhost:8000/set', { key, value })
      .then(function (response) {
        console.log('Response:', response.data);
        setCacheAlert(response.data.STATUS)
        setTimeout(() =>{
          setCacheAlert("")
        }, 2000)
       
      })
      .catch(function (error) {
        console.error('Error:', error);
        
      });
  }
  const handleGetValue = () =>{
    console.log(getValue)
    axios.get(`http://localhost:8000/get/${getValue}`)
    .then(function (response) {
      console.log('GET Response:', response.data);
      if (response.data.STATUS == "DATA FOUND") {
        SetGetValueAlert(response.data.DATA)
        setTimeout(() =>{
          SetGetValueAlert("")
        }, 2000)
      } else {
        SetError(true)
        setTimeout(() =>{
          SetError(false)
        }, 2000)
      }

      // Handle the GET response data as needed
    })
    .catch(function (error) {
      console.error('GET Error:', error);
      // Handle GET errors
    });
  }
  return (
    <div className="App" >
      <h1>LRU cacheing system</h1>
      <Paper elevation={3} sx={{ width: "800px", height: "300px", margin: "0 auto", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column", gap: "10px", backgroundColor: "#F5EEE6", position: "relative", padding: "20px", marginBottom: "20px" }}>
        <h2>SET cache</h2>
        <TextField id="outlined-basic" label="KEY" variant="outlined" onChange={(e) => { setKey(e.target.value) }} />
        <TextField id="outlined-basic" label="VALUE" variant="outlined" onChange={(e) => { setValue(e.target.value) }} />
        {key !== "" && value !== "" ? <Button variant="contained" onClick={handleSendButton}>SEND</Button> : <Button variant="contained" disabled>SEND</Button>}
        {cacheAlert !== "" && <Alert severity="success" sx={{ position: "absolute", bottom: "10px", left: "50%", transform: "translateX(-50%)" }}>Cache inserted successfully</Alert>}
      </Paper>

      <Paper elevation={3} sx={{ width: "800px", height: "300px", margin: "0 auto", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column", gap: "10px", backgroundColor: "#F5EEE6", position: "relative", padding: "20px" }}>
        <h2>GET cache</h2>
        <TextField id="outlined-basic" label="KEY" variant="outlined" onChange={(e)=>{setGetValue(e.target.value)}}/>
        {getValue !== "" ? <Button variant="contained" onClick={handleGetValue}>GET VALUE</Button>: <Button variant="contained" disabled>GET VALUE</Button>}
        {getValueAlert !== "" && <Alert severity="success" sx={{ position: "absolute", bottom: "10px", left: "50%", transform: "translateX(-50%)" }}>VALUE: {getValueAlert}</Alert>}
        {error == true &&  <Alert severity="error">Cache not found or expired!!!</Alert>}
      </Paper>


      
      
    </div>
  );
}

export default App;
