import React from 'react';
import ReactDOM from 'react-dom/client';
import "./css/app.css"
import "./main.js"
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import SearchIcon from '@mui/icons-material/Search';
import Grid from '@mui/material/Grid';
import Container from '@mui/material/Container';
import TextField from '@mui/material/TextField';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';
import ImageIcon from '@mui/icons-material/Image';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <AppBar position="static" elevation={2}>
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          OnlineChat
        </Typography>
      </Toolbar>
    </AppBar>

    <Container maxWidth="lg">
      <Paper className='chatPaper' elevation={1}>
        <Grid className='grid' container spacing={0}>
          <Grid lg={3} sm={4} sx={{ display: { xs: 'none', sm: 'block' } }}>
            <AppBar className='appbar' position="static" elevation={0}>
              <Toolbar>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>群组成员</Typography>
                <IconButton
                  size="large"
                  edge="start"
                  color="inherit"
                  aria-label="menu"
                // sx={{ mr: 2 }}
                >
                  <SearchIcon />
                </IconButton>
              </Toolbar>
            </AppBar>
            <List className='userList' sx={{ width: '100%', height: '100%' }}>
              <ListItem>
                <ListItemAvatar>
                  <Avatar>
                    <ImageIcon />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText primary="Offline" secondary="在线状态..." />
              </ListItem>
              <TextField id="name" label="名字" value={'000'} variant="filled" />
            </List>
          </Grid>
          <Grid lg={9} xs={12} sm={8}>
            <AppBar className='chat' position="static" elevation={0}>
              <Toolbar>
                <IconButton
                  size="large"
                  edge="start"
                  color="inherit"
                  aria-label="menu"
                  sx={{ display: { xs: 'flex', sm: 'none' }, mr: 2 }}
                >
                  <MenuIcon />
                </IconButton>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                  在线聊天
                </Typography>
              </Toolbar>
            </AppBar>
            <div id="log"></div>
            <form id="form">
              <TextField id="msg" label="消息" variant="filled" />
              <IconButton
                size="large"
                edge="start"
                color="inherit"
                aria-label="menu"
                type="submit"
                sx={{ mr: 2 }}
              >
                <MenuIcon />
              </IconButton>
            </form>
          </Grid>
        </Grid>
      </Paper>
    </Container>
  </React.StrictMode>
);
