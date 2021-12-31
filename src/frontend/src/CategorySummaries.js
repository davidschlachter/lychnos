import React from 'react';
import './CategorySummaries.css';
import AmountLeft from './AmountLeft';
import FillBar from './FillBar.js'
import Spinner from './Spinner.js'
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { Link } from 'react-router-dom';

class CategorySummaries extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            summaries: [],
            budgets: []
        };
    }

    componentDidMount() {
        fetch("/api/budgets/")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        budgets: result
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
        fetch("/api/reports/categorysummary/?budget=1")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        summaries: result
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
    }

    render() {
        const { error, isLoaded, summaries, budgets } = this.state;
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <Spinner />;
        } else {
            let timeSpent = 0
            if (budgets.length > 1) {
                let start = new Date(budgets[0].start)
                let end = new Date(budgets[0].end)
                let now = new Date()
                timeSpent = (Math.abs(now - start) / Math.abs(end - start)) * 100
            }
            return (
                <Box sx={{ p: 2, mb: 6 }}>
                    <TableContainer component={Paper}>
                        <Table style={{ "width": "100%" }}>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Category</TableCell>
                                    <TableCell>Progress</TableCell>
                                    <TableCell>Left&nbsp;per month</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {summaries.map(item => (
                                    <TableRow key={item.category_budget_id} to={"/categorydetail/" + item.category_budget_id} component={Link}>
                                        <TableCell>{item.name}</TableCell>
                                        <TableCell width="99%"><FillBar amount={item.amount} sum={item.sum} now={timeSpent}></FillBar></TableCell>
                                        <TableCell style={{
                                            'fontFamily': "monospace",
                                            'textAlign': "right",
                                            'fontSize': "110%",
                                            'fontWeight': "600"
                                        }}><AmountLeft amount={item.amount} sum={item.sum} timeSpent={timeSpent}></AmountLeft></TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Box>
            );
        }
    }
}

export default CategorySummaries;
