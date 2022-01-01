function AmountLeft(props) {
    let timeSpent;
    if (props.timeSpent > 100) {
        timeSpent = 100;
    } else {
        timeSpent = props.timeSpent
    }

    let left;
    if (timeSpent == 100) {
        left = Math.round(props.amount - props.sum)
    } else {
        left = Math.round((props.amount - props.sum) / Math.ceil((1 - timeSpent / 100) * 12))
    }
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
