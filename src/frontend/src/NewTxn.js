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
    const [sourceAccount, setSourceAccount] = React.useState('');
    const [destinationAccount, setDestinationAccount] = React.useState('');
    const [description, setDescription] = React.useState('');
    const [amount, setAmount] = React.useState('');

    const accountsAreSame = sourceAccount && destinationAccount && sourceAccount === destinationAccount;
    const isFormValid = description.trim() && amount.trim() && sourceAccount.trim() && destinationAccount.trim() && !accountsAreSame;

    function handleClick(event) {
        if (!isFormValid) {
            event.preventDefault();
            return false;
        }
        setSubmitted(true);
        return true;
    }

    let now = new Date()

    return (
        <>
            <Header back_visibility="hidden" title="New transaction"></Header>
            <Box sx={{ p: 2, pb: 8 }} component="form" action="/api/transactions/" method="POST" onSubmit={handleClick} >
                <Stack direction="column">
                    <TextField
                        required
                        name="description"
                        id="description"
                        label="Description"
                        variant="outlined"
                        align="center"
                        margin="normal"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                    />
                    <AccountsInput
                        name="source_name"
                        label="Source account"
                        value={sourceAccount}
                        onInputChange={(event, newValue) => setSourceAccount(newValue)}
                    />
                    <AccountsInput
                        name="destination_name"
                        label="Destination account"
                        value={destinationAccount}
                        onInputChange={(event, newValue) => setDestinationAccount(newValue)}
                    />
                    <TextField
                        required
                        id="date"
                        name="date"
                        type="date"
                        variant="outlined"
                        align="center"
                        margin="normal"
                        // en-CA is different between iOS Safari and every other
                        // browser! However, fr-CA seems to always be
                        // yyyy-mm-dd.
                        defaultValue={now.toLocaleDateString('fr-CA')}
                    />
                    <CategoriesInput name="category_name" label="Category" />
                    <TextField
                        required
                        id="amount"
                        name="amount"
                        label="Amount"
                        variant="outlined"
                        align="center"
                        margin="normal"
                        autoComplete="off"
                        value={amount}
                        onChange={(e) => setAmount(e.target.value)}
                        slotProps={{
                            htmlInput: { inputMode: 'decimal', pattern: '[0-9.,]*' }
                        }}
                    />
                    <LoadingButton
                        variant="contained"
                        name="submitButton"
                        type="submit"
                        loading={submitted}
                        disabled={submitted || !isFormValid}
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