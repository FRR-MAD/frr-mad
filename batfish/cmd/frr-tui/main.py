import os
import shutil
from rich.console import Console
from rich.layout import Layout
from rich.panel import Panel
from rich.table import Table
from rich.prompt import Prompt
from internal.utils.batfish_utils import (
    setup_batfish,
    list_snapshots,
    detect_bgp_anomalies,
    detect_ospf_anomalies,
    ensure_snapshot_structure,
)

console = Console()

def create_layout():
    """Create the TUI layout."""
    layout = Layout()
    layout.split(
        Layout(name="header", size=3),
        Layout(name="main", ratio=1),
        Layout(name="footer", size=3)
    )
    layout["header"].update(Panel("FRR TUI - Network Anomaly Monitoring", style="bold blue"))
    layout["footer"].update(Panel("Press Q to quit", style="bold red"))
    return layout

def display_results(results):
    """Display analysis results in a rich table."""
    for anomaly_type, result in results.items():
        if result.empty:
            console.print(f"No {anomaly_type.replace('_', ' ')} found.", style="green")
            continue
            
        table = Table(title=f"{anomaly_type.replace('_', ' ').title()}")

        # Add columns
        for column in result.columns:
            table.add_column(column)

        # Add rows
        for _, row in result.iterrows():
            table.add_row(*[str(row[col]) for col in result.columns])

        console.print(table)

def main():
    """Main TUI function."""
    # Setup basic directories if not present
    os.makedirs("configs/snapshots", exist_ok=True)

    try:
        # Initialize Batfish
        bf = setup_batfish()  # Get initialized session

        while True:
            layout = create_layout()
            snapshots = list_snapshots()

            # Create a default snapshot if none exists
            if not snapshots:
                console.print("No snapshots found. Creating default snapshot directory.", style="yellow")
                default_snapshot_path = "configs/snapshots/default"
                ensure_snapshot_structure(default_snapshot_path)
                snapshots = ["default"]
                console.print("Please add router configuration files to:", style="yellow")
                console.print(f"{os.path.abspath(os.path.join(default_snapshot_path, 'configs'))}", style="bold yellow")
                console.print("Then restart the application.", style="yellow")
                input("Press Enter to exit...")
                break

            # Let the user select a snapshot
            console.print("Available snapshots:", style="bold blue")
            for idx, snapshot in enumerate(snapshots, 1):
                console.print(f"{idx}. {snapshot}")
            console.print(f"Q. Quit")
            
            choice = Prompt.ask("Select a snapshot", choices=[str(i) for i in range(1, len(snapshots)+1)] + ["Q"])
            if choice == "Q":
                break
                
            snapshot_name = snapshots[int(choice)-1]

            # Initialize the selected snapshot with proper path structure
            snapshot_path = os.path.abspath(os.path.join("configs/snapshots", snapshot_name))
            ensure_snapshot_structure(snapshot_path)  # Make sure it has the correct structure
            
            try:
                bf.init_snapshot(snapshot_path, name=snapshot_name, overwrite=True)
                console.print(f"Initialized snapshot: {snapshot_name}", style="green")
            except Exception as e:
                console.print(f"Error initializing snapshot: {str(e)}", style="bold red")
                input("Press Enter to continue...")
                continue

            # Let the user select a protocol
            protocol = Prompt.ask("Select protocol", choices=["bgp", "ospf", "Q"])
            if protocol == "Q":
                break

            # Run pre-screening and detect anomalies
            try:
                if protocol == "bgp":
                    anomalies = detect_bgp_anomalies()
                elif protocol == "ospf":
                    anomalies = detect_ospf_anomalies()
                    
                # Display the results
                display_results(anomalies)
            except Exception as e:
                console.print(f"Error analyzing {protocol.upper()} configurations: {str(e)}", style="bold red")
            
            # Wait for user input before continuing
            if Prompt.ask("Press Enter to continue or Q to quit") == "Q":
                break
    except Exception as e:
        console.print(f"Error initializing Batfish: {str(e)}", style="bold red")
        console.print("Please ensure you have valid router configurations in at least one snapshot.", style="yellow")
        console.print("You can add .cfg or .conf files to the 'configs/snapshots/my_snapshot/configs' directory.", style="yellow")
        if Prompt.ask("Would you like to continue without Batfish initialized?", choices=["y", "n"]) == "n":
            return

if __name__ == "__main__":
    main()