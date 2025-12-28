import * as React from 'react';
import useFetch from "./useFetch";
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';

export default function CategoriesInput(props) {
    const { response, error, loading } = useFetch(
        "/api/categories/",
        {
            query: {},
        }
    );

    if (loading) {
        return (
            <TextField
                disabled
                label={props.label}
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
                name={props.name}
                margin="normal"
                autoComplete="off"
                variant="outlined"
            />
        );
    }
    return (
        <Autocomplete
            disablePortal
            id={props.name}
            openOnFocus
            disableClearable={false}
            options={response.filter(c => c.ignore == false).map(c => (c.name)).sort()}
            sx={{ width: 300 }}
            autoComplete="off"
            renderInput={(params) => <TextField
                {...params}
                label={props.label}
                name={props.name}
                margin="normal"
                autoComplete="off"
                variant="outlined"
            />}
        />
    );
}