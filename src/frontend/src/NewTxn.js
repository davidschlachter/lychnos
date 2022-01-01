import * as React from 'react';
import AccountsInput from "./AccountsInput.js";
import CategoriesInput from "./CategoriesInput.js";
import Header from './Header.js';
import Box from '@mui/material/Box';
import LoadingButton from '@mui/lab/LoadingButton';
import SaveIcon from '@mui/icons-material/Save';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';

export default function NewTxn() {
    const [submitted, setSubmitted] = React.useState(false);
    function handleClick() {
        setSubmitted(true);
        return true;
    }

    let now = new Date()

    return (
        <>
            <Header back_location="/" title="New transaction"></Header>
            <Box sx={{ p: 2, mb: 6 }} component="form" action="/api/transactions/" method="POST" onSubmit={handleClick} >
                <Stack direction="column">
                    <TextField
                        required
                        name="description"
                        id="description"
                        label="Description"
                        variant="outlined"
                        align="center"
                        margin="normal"
                    />
                    <AccountsInput name="source_name" label="Source account" />
                    <AccountsInput name="destination_name" label="Destination account" />
                    <TextField
                        required
                        id="date"
                        name="date"
                        type="date"
                        variant="outlined"
                        align="center"
                        margin="normal"
                        defaultValue={now.toLocaleDateString('en-CA')}
                    />
                    <TextField
                        required
                        id="amount"
                        name="amount"
                        inputProps={{ inputMode: 'numeric', pattern: '[0-9.]*' }}
                        label="Amount"
                        variant="outlined"
                        align="center"
                        margin="normal"
                        autoComplete="off"
                    />
                    <CategoriesInput name="category_name" label="Category" />
                    <LoadingButton
                        variant="contained"
                        name="submitButton"
                        type="submit"
                        loading={submitted}
                        disabled={submitted}
                        startIcon={<SaveIcon />}
                        loadingPosition="start"
                        sx={{ m: 1 }}
                    >
                        Save
                    </LoadingButton>
                </Stack>
            </Box>
        </>
    );
}