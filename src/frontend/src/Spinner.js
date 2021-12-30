import * as React from 'react';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';

export default function Spinner() {
    return (
        <Box sx={{ mx: 'auto', textAlign: 'center', p: '3em' }}>
            <CircularProgress />
        </Box>
    );
}