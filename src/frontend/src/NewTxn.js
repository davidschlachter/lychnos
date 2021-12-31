import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Header from './Header.js';

export default function NewTxn() {
    return (
        <>
            <Header back_location="/" title="New transaction"></Header>
            <Box sx={{ p: 2 }} component="form">
                <Stack direction="column">
                    <TextField id="description" label="Description" variant="outlined" align="center" margin="normal" />
                    <TextField id="from" label="Source account" variant="outlined" align="center" margin="normal" />
                    <TextField id="to" label="Destination account" variant="outlined" align="center" margin="normal" />
                    <TextField id="date" type="date" variant="outlined" align="center" margin="normal" />
                    <TextField id="amount" type="number" min="0.00" step="0.01" label="Amount" variant="outlined" align="center" margin="normal" />
                    <TextField id="category" label="Category" variant="outlined" align="center" margin="normal" />
                    <Button variant="contained">Save</Button>
                </Stack>
            </Box>
        </>
    );
}