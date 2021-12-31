export default function FriendlyDate(props) {
    let options = { month: 'short', day: 'numeric' };
    let date = new Date(props.date)
    let dateStr = date.toLocaleDateString("en-US", options)
    return (
        <span>{dateStr}</span>
    );
}