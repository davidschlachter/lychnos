import TableCell from '@mui/material/TableCell';

function AmountLeft(props) {
    let left = Math.round((props.amount - props.sum) / Math.ceil((1 - props.timeSpent / 100) * 12))
    let color = 'black'
    let sign = 1
    if (Math.sign(left) !== Math.sign(props.amount)) {
        if (Math.abs(left) > 10) { color = 'red' }
        sign = -1
    }

    return (
        <span style={{
            "color": color
        }}>{sign * Math.abs(left)}</span>
    );
}

export default AmountLeft;
