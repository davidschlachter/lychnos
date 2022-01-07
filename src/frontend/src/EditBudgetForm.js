import React from 'react';
import EditCategoryRow from './EditCategoryRow.js';
import Box from '@mui/material/Box';
import LoadingButton from '@mui/lab/LoadingButton';
import Paper from '@mui/material/Paper';
import SaveIcon from '@mui/icons-material/Save';
import Stack from '@mui/material/Stack';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';

export default function EditBudget(props) {
    const [error, setError] = React.useState(false);
    const [submitted, setSubmitted] = React.useState(false);
    const [budget, setBudget] = React.useState(props.seed);
    const [sums, setSums] = React.useState({ income: sumOnSign(1), expenses: sumOnSign(-1), net: sumAll() });

    const categories = [...props.categories].sort((a, b) => (a.name > b.name ? 1 : -1));
    const categorybudgets = props.categorybudgets;

    function updateBudget(category, amount) {
        let b = budget;
        b[category] = amount;
        setBudget(b);
        setSums({ income: sumOnSign(1), expenses: sumOnSign(-1), net: sumAll() });
    }
    function submitForm() {
        setSubmitted(true)
        let req = [];
        for (const [key, value] of Object.entries(budget)) {
            req.push({ category: parseInt(key), amount: value })
        }
        (async () => {
            const rawResponse = await fetch('/api/categorybudgets/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(req)
            });
            const status = await rawResponse.status;
            if (status === 201) {
                window.location.href = "/app/";
            } else {
                setError(rawResponse.body)
            }
        })();
    }
    function sumOnSign(sign) {
        let sum = 0.0;
        for (const [id, amount] of Object.entries(budget)) {
            if (Math.sign(amount) == sign) {
                sum = parseFloat(sum) + parseFloat(amount);
            }
        }
        return sum;
    }
    function sumAll() {
        let sum = 0;
        for (const [id, amount] of Object.entries(budget)) {
            sum = parseFloat(sum) + parseFloat(amount);
        }
        return sum;
    }

    if (error) {
        return (
            <p>Error: {error}</p>
        );
    }

    return (
        <>
            <Box sx={{ p: 2, pb: 8 }}>
                <Typography variant="subtitle1" component="div" align="center" gutterBottom>
                    Income: ${sums.income}, Expenses: ${sums.expenses}<br />
                    Net: ${sums.net}
                </Typography>
                <TableContainer component={Paper}>
                    <Table style={{ "width": "100%" }}>
                        <TableHead>
                            <TableRow>
                                <TableCell>Category</TableCell>
                                <TableCell>Budget</TableCell>
                                <TableCell>Monthly</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {categories.map(item => (
                                <EditCategoryRow id={item.id} name={item.name} amount={getAmount(item.id, categorybudgets)} updateFunc={updateBudget} />
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                <Stack direction="column" justifyContent="center">
                    <LoadingButton
                        variant="contained"
                        name="submitButton"
                        size="large"
                        align="center"
                        loading={submitted}
                        disabled={submitted}
                        startIcon={<SaveIcon />}
                        loadingPosition="start"
                        sx={{ m: 1 }}
                        onClick={submitForm}
                    >
                        Save
                    </LoadingButton>
                </Stack>
            </Box>
        </>
    );
}

function getAmount(categoryID, categorybudgets) {
    for (const cb of categorybudgets) {
        if (cb.category === categoryID) {
            return cb.amount;
        }
    }
    return 0;
}