import React from 'react';
import FillBar from './FillBar.js'
import CategoryHeader from './CategoryHeader.js'
import AmountLeft from './AmountLeft';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';

class CategoryDetail extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            details: [],
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
        fetch("/api/reports/categorysummary/" + String(this.props.categoryId))
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        details: result
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
        const { error, isLoaded, details, budgets } = this.state;
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            let timeSpent = 0
            if (budgets.length > 1) {
                let start = new Date(budgets[0].start)
                let end = new Date(budgets[0].end)
                let now = new Date()
                timeSpent = (Math.abs(now - start) / Math.abs(end - start)) * 100
            }
            let totalSpent = 0
            details[0].totals.map(item => (
                totalSpent += Math.round(parseFloat(item.earned) + parseFloat(item.spent))
            ));
            return (
                <Paper>
                    <CategoryHeader category_name={details[0].name}></CategoryHeader>
                    <Typography variant="h6" component="div" align="center" gutterBottom>
                        Budgeted: {details[0].amount}, Actual: {totalSpent}<br />
                        Left per month: <AmountLeft amount={details[0].amount} sum={totalSpent} timeSpent={timeSpent} />
                    </Typography>
                    <div style={{ "width": "100%;" }}>
                        <FillBar amount={details[0].amount} sum={details[0].sum} now={timeSpent}></FillBar>
                    </div>
                    <TableContainer component={Paper}>
                        <Table style={{ "width": "100%", "margin-bottom": "3.5em" }}>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Month</TableCell>
                                    <TableCell>Sum</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {details[0].totals.map(item => (
                                    <TableRow key={item.start}>
                                        <TableCell>{Intl.DateTimeFormat('en', { month: 'long' }).format(new Date(item.start.substring(5, 7)))}</TableCell>
                                        <TableCell>{Math.round(parseFloat(item.earned) + parseFloat(item.spent))}</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Paper>
            );
        }
    }
}

export default CategoryDetail;