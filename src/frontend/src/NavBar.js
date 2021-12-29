import * as React from 'react';
import Paper from '@mui/material/Paper';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import AlignHorizontalLeftIcon from '@mui/icons-material/AlignHorizontalLeft';
import AddBoxIcon from '@mui/icons-material/AddBox';
import ListAltIcon from '@mui/icons-material/ListAlt';

export default function NavBar() {
    const [value, setValue] = React.useState(0);

    return (
        <Paper sx={{ position: 'fixed', bottom: 0, left: 0, right: 0 }} elevation={3}>
            <BottomNavigation
                showLabels
                value={value}
                onChange={(event, newValue) => {
                    setValue(newValue);
                }}
            >
                <BottomNavigationAction label="Summary" icon={<AlignHorizontalLeftIcon />} />
                <BottomNavigationAction label="New Txn" icon={<AddBoxIcon />} />
                <BottomNavigationAction label="Txns" icon={<ListAltIcon />} />
            </BottomNavigation>
        </Paper>
    );
}