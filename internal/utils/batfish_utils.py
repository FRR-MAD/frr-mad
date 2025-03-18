import os
import shutil
from pybatfish.client.session import Session

# Initialize Batfish session globally
bf = Session(host="localhost")  # Set Batfish server host

def setup_batfish():
    """Initialize the Pybatfish session."""
    bf.set_network("my_network")  # Set a network name
    
    # Ensure the snapshot directory exists with proper structure
    snapshot_path = "configs/snapshots/my_snapshot"
    ensure_snapshot_structure(snapshot_path)
    
    # Check if the configs directory is empty
    configs_dir = os.path.join(snapshot_path, "configs")
    if not os.listdir(configs_dir):
        print("Config directory is empty. Copying sample configurations...")
        
        # Find a sample configuration from another snapshot
        sample_found = False
        for sample_snapshot in ["bgp_duplicated", "ospf_duplicated"]:
            sample_path = f"configs/snapshots/{sample_snapshot}/configs"
            
            # Try with "configs" directory first (if you renamed them)
            if os.path.exists(sample_path) and os.listdir(sample_path):
                sample_found = True
            else:
                # Fall back to check "config" directory (if you haven't renamed them yet)
                sample_path = f"configs/snapshots/{sample_snapshot}/config"
                if os.path.exists(sample_path) and os.listdir(sample_path):
                    sample_found = True
            
            if sample_found:
                # Copy all configuration files from the sample
                for filename in os.listdir(sample_path):
                    if filename.endswith(('.cfg', '.conf')) or not filename.endswith(('.txt', '.md')):
                        source_path = os.path.join(sample_path, filename)
                        dest_path = os.path.join(configs_dir, filename)
                        if os.path.isfile(source_path):
                            shutil.copy2(source_path, dest_path)
                            print(f"Copied {filename} to {dest_path}")
                break
        
        if not sample_found:
            print("Warning: Could not find sample configurations. Creating a minimal configuration...")
            # Create a minimal configuration file
            with open(os.path.join(configs_dir, "minimal.cfg"), "w") as f:
                f.write("hostname minimal-router\n")
                f.write("!\n")
                f.write("interface Loopback0\n")
                f.write(" ip address 192.168.1.1/32\n")
                f.write("!\n")
            print("Created minimal configuration file.")
    
    # Initialize the snapshot
    bf.init_snapshot(snapshot_path, name="my_snapshot", overwrite=True)
    print("Batfish session initialized.")
    return bf  # Return the session object for use in main

def ensure_snapshot_structure(snapshot_path):
    """
    Ensure the snapshot directory has the expected structure.
    Batfish expects at least one of: hosts, configs, aws_configs, or sonic_configs directories.
    """
    # Create the main snapshot directory if it doesn't exist
    os.makedirs(snapshot_path, exist_ok=True)
    
    # Check for legacy "config" directory and rename to "configs" if found
    legacy_config_dir = os.path.join(snapshot_path, "config")
    configs_dir = os.path.join(snapshot_path, "configs")
    
    if os.path.exists(legacy_config_dir) and not os.path.exists(configs_dir):
        # Rename legacy directory to the new format
        os.rename(legacy_config_dir, configs_dir)
        print(f"Renamed legacy 'config' to 'configs' in {snapshot_path}")
    else:
        # Create the configs subdirectory if it doesn't exist
        os.makedirs(configs_dir, exist_ok=True)
    
    # If there are any .cfg or .conf files in the snapshot directory, move them to configs
    for filename in os.listdir(snapshot_path):
        if filename.endswith(('.cfg', '.conf')):
            source_path = os.path.join(snapshot_path, filename)
            dest_path = os.path.join(configs_dir, filename)
            if os.path.isfile(source_path) and not os.path.exists(dest_path):
                shutil.copy2(source_path, dest_path)
    
    print(f"Snapshot structure created at {snapshot_path}")

def list_snapshots():
    """List available snapshots in the configs/snapshots/ directory."""
    snapshots_dir = "configs/snapshots"
    if not os.path.exists(snapshots_dir):
        os.makedirs(snapshots_dir, exist_ok=True)
    return [name for name in os.listdir(snapshots_dir) if os.path.isdir(os.path.join(snapshots_dir, name))]

def pre_screen_bgp():
    """Run basic BGP checks."""
    return {
        "missing_neighbors": bf.q.bgpPeerConfiguration().answer().frame(),
        "missing_networks": bf.q.bgpProcessConfiguration().answer().frame(),
    }

def pre_screen_ospf():
    """Run basic OSPF checks."""
    return {
        "missing_interfaces": bf.q.ospfInterfaceConfiguration().answer().frame(),
        "area_config": bf.q.ospfAreaConfiguration().answer().frame(),
    }

def detect_bgp_anomalies():
    """Detect BGP anomalies based on pre-screening results."""
    pre_screen_results = pre_screen_bgp()
    anomalies = {}

    # Basic BGP session status
    anomalies["bgp_sessions"] = bf.q.bgpSessionStatus().answer().frame()
    
    # Routing policy issues
    anomalies["undefined_references"] = bf.q.undefinedReferences().answer().frame()
    
    # Route advertisement issues
    if not pre_screen_results["missing_networks"].empty:
        anomalies["unused_structures"] = bf.q.unusedStructures().answer().frame()
        anomalies["bgp_routes"] = bf.q.routes(protocols="BGP").answer().frame()
        
        # Check for BGP route reflection configuration
        anomalies["bgp_route_reflection"] = bf.q.bgpRib().answer().frame()

    return anomalies

def detect_ospf_anomalies():
    """Detect OSPF anomalies based on pre-screening results."""
    pre_screen_results = pre_screen_ospf()
    anomalies = {}

    # OSPF session status
    anomalies["ospf_neighbors"] = bf.q.ospfSessionCompatibility().answer().frame()
    
    # OSPF area configuration issues
    anomalies["area_configuration"] = bf.q.ospfAreaConfiguration().answer().frame()
    
    # Route advertisement issues
    if not pre_screen_results["area_config"].empty:
        anomalies["ospf_routes"] = bf.q.routes(protocols="OSPF").answer().frame()
        anomalies["interface_properties"] = bf.q.interfaceProperties().answer().frame()

    return anomalies