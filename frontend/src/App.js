import axios from 'axios';
import React, { useState, useEffect } from 'react';
import { Box, Button, Paper, TextField, Alert, Select, MenuItem } from '@mui/material';
import './App.css';

function App() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [getValue, setGetValue] = useState("");
  const [cacheAlert, setCacheAlert] = useState("");
  const [getValueAlert, SetGetValueAlert] = useState("");
  const [error, SetError] = useState(false);
  const [expiration, setExpiration] = useState(5); // default to 5 seconds
  const [cacheData, setCacheData] = useState({});

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8000/ws');
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setCacheData(data);
    };
    return () => socket.close();
  }, []);

  const handleSendButton = () => {
    console.log(key, value, expiration);

    axios.post('http://localhost:8000/set', { key, value, expiration })
      .then(function (response) {
        console.log('Response:', response.data);
        setCacheAlert(response.data.STATUS);
        setTimeout(() => {
          setCacheAlert("");
        }, 2000);
      })
      .catch(function (error) {
        console.error('Error:', error);
      });
  };

  const handleGetValue = () => {
    console.log(getValue);
    axios.get(`http://localhost:8000/get/${getValue}`)
      .then(function (response) {
        console.log('GET Response:', response.data);
        if (response.data.STATUS === "DATA FOUND") {
          SetGetValueAlert(response.data.DATA);
          setTimeout(() => {
            SetGetValueAlert("");
          }, 2000);
        } else {
          SetError(true);
          setTimeout(() => {
            SetError(false);
          }, 2000);
        }
      })
      .catch(function (error) {
        console.error('GET Error:', error);
      });
  };

  return (
    <div className="App">
      <h1>LRU Caching System</h1>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%', padding: '20px' }}>
        <Box sx={{ width: '60%' }}>
          <Paper elevation={3} sx={{ width: "100%", height: "350px", margin: "0 auto", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column", gap: "10px", backgroundColor: "#F5EEE6", position: "relative", padding: "20px", marginBottom: "20px" }}>
            <h2>SET Cache</h2>
            <TextField id="outlined-basic" label="KEY" variant="outlined" onChange={(e) => { setKey(e.target.value) }} />
            <TextField id="outlined-basic" label="VALUE" variant="outlined" onChange={(e) => { setValue(e.target.value) }} />
            <Select
              value={expiration}
              onChange={(e) => setExpiration(e.target.value)}
              displayEmpty
              inputProps={{ 'aria-label': 'Without label' }}
              sx={{ width: "200px" }}
            >
              {[...Array(9).keys()].map(i => (
                <MenuItem key={i + 1} value={i + 1}>{i + 1} Seconds</MenuItem>
              ))}
            </Select>
            {key !== "" && value !== "" ? <Button variant="contained" onClick={handleSendButton}>SEND</Button> : <Button variant="contained" disabled>SEND</Button>}
            {cacheAlert !== "" && <Alert severity="success" sx={{ position: "absolute", bottom: "10px", left: "50%", transform: "translateX(-50%)" }}>Cache inserted successfully</Alert>}
          </Paper>

          <Paper elevation={3} sx={{ width: "100%", height: "300px", margin: "0 auto", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column", gap: "10px", backgroundColor: "#F5EEE6", position: "relative", padding: "20px" }}>
            <h2>GET Cache</h2>
            <TextField id="outlined-basic" label="KEY" variant="outlined" onChange={(e) => { setGetValue(e.target.value) }} />
            {getValue !== "" ? <Button variant="contained" onClick={handleGetValue}>GET VALUE</Button> : <Button variant="contained" disabled>GET VALUE</Button>}
            {getValueAlert !== "" && <Alert severity="success" sx={{ position: "absolute", bottom: "10px", left: "50%", transform: "translateX(-50%)" }}>VALUE: {getValueAlert}</Alert>}
            {error === true && <Alert severity="error">Cache not found or expired!!!</Alert>}
          </Paper>
        </Box>

        <Box sx={{ width: '35%' }}>
          <Paper elevation={3} sx={{ width: "100%", height: "400px", margin: "0 auto", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column", gap: "10px", backgroundColor: "#F5EEE6", position: "relative", padding: "20px" }}>
            <h2>Current Cache</h2>
            <Box sx={{ width: "100%", maxHeight: "300px", overflow: "auto" }}>
              {Object.keys(cacheData).length === 0 ? (
                <p>No cache data available.</p>
              ) : (
                <ul>
                  {Object.entries(cacheData).map(([key, value]) => (
                    <li key={key}><strong>{key}:</strong> {value}</li>
                  ))}
                </ul>
              )}
            </Box>
          </Paper>
        </Box>
      </Box>
    </div>
  );
}

export default App;
