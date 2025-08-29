import * as React from 'react';
import useFetch from "./useFetch";
import Autocomplete from '@mui/material/Autocomplete';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';

export default function AccountsInput(props) {
    const { response, error, loading } = useFetch(
        "/api/accounts/?type=asset&type=revenue&type=expense",
        {
            query: {},
        }
    );

    if (loading) {
        return (
            <TextField
                disabled
                label={props.label}
                required
                name={props.name}
                margin="normal"
                variant="outlined"
            />
        );
    }
    if (error) {
        return (
            <TextField
                label={props.label}
                required
                name={props.name}
                margin="normal"
                variant="outlined"
            />
        );
    }

    // Tidy up the autocomplete arrays
    let raw_options = []
    for (const e of response) {
        let item = { label: e.attributes.name }
        if (e.attributes.type === "asset") {
            item.display_string = e.attributes.name + " ($" + e.attributes.current_balance + ")"
        } else {
            item.display_string = item.label
        }
        raw_options.push(item)
    }
    // Remove any duplicates
    let account_options = []
    let buffer = new Set()
    for (const e of raw_options) {
        if (!buffer.has(e.label)) {
            buffer.add(e.label)
            account_options.push(e)
        }
    }

    return (
        <Autocomplete
            freeSolo
            id={props.name}
            openOnFocus
            options={account_options}
            value={props.value || ''}
            onInputChange={props.onInputChange}
            getOptionLabel={(option) => option.label}
            renderOption={(props, option) => (
                <Box component="li" {...props}>
                    {option.display_string}
                </Box>
            )}
            renderInput={(params) => (
                <TextField
                    {...params}
                    label={props.label}
                    required
                    name={props.name}
                    margin="normal"
                    variant="outlined"
                    autoComplete="off"
                    slotProps={{
                        input: {
                            ...params.InputProps,
                            type: 'search',
                        }
                    }}
                />
            )}
        />
    );

}