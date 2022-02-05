import './FillBar.css';

const green = '#a2ff00';
const yellow = '#ffdd00';
const red = '#ff5e00';

function FillBar(props) {
    let width = props.sum / props.amount * 100
    let color = green;
    if (props.sum < 0 && props.amount > 0) {
        width = 100;          // Bar width is meaningless if sign(sum) != sign(amount)
        color = red;          // Expected net positive, currently net negative.

    } else if (props.sum > 0 && props.amount < 0) {
        width = 0;            // Expected net negative, currently net positive.
    } else if (width > 100) { // Category is beyond budget (but has expected sign)
        width = 100;
        if (props.sum > 0) {  // Good for income
            color = green;
        } else {              // Bad for expenses
            color = red;
        }
    } else if (props.amount < 0 && Math.abs(props.sum / props.amount) > ((1 / 12) + (props.now / 100))) {
        color = yellow;       // Expense target is running beyond target
    } else if (props.amount > 0 && Math.abs(props.sum / props.amount) < ((props.now / 100) - (3 / 12))) {
        color = yellow;       // Income category is running behind target (more generous here, accommodate seasonal work)
    }
    let style = {
        'width': String(width) + "%",
        'backgroundColor': color
    }

    let nowStyle = {
        'left': String(props.now) + "%",
        'width': "2px",
        'height': "2em",
        'position': "relative",
        'backgroundColor': 'black',
        'top': "-2em"
    }

    return (
        <div className="fillBar">
            <div className="filled" style={style}></div>
            <div className="nowLine" style={nowStyle}></div>
        </div>
    );
}

export default FillBar;
