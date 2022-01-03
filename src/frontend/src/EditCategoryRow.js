import React from 'react';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import TextField from '@mui/material/TextField';

export default function EditCategoryRow(props) {
    const [amount, setAmount] = React.useState(props.amount);
    function handleInput(event) {
        setAmount(event.target.value)
        props.updateFunc(props.id, event.target.value)
    }
    const [disabled, toggleDisabled] = React.useState(props.amount === 0);
    function handleToggle() {
        if (disabled) {
            toggleDisabled(false)
        } else {
            toggleDisabled(true)
        }

    }

    return (

        <TableRow key={props.id}>
            <TableCell><input
                name={props.name}
                type="checkbox"
                checked={!disabled}
                onClick={handleToggle}
            /> {props.name}</TableCell>
            <TableCell><TextField
                disabled={disabled}
                required
                id={props.id + "amount"}
                name="amount"
                inputProps={{ inputMode: 'numeric', pattern: '[0-9]*' }}
                label="Amount"
                variant="outlined"
                align="center"
                margin="normal"
                autoComplete="off"
                defaultValue={props.amount}
                onInput={handleInput}
            /></TableCell>
            <TableCell>{Math.round(amount / 12)}</TableCell>
        </TableRow>
    );
}