# Helpers ----------------------------------------------------------------------

def pretty_version(version):
    return "{0}.{1}".format(version["major"], version["minor"])

def no_empty_log_in_change(change):
    version = pretty_version(change["version"])
    log = change["log"]

    if len(log) == 0:
        return "Logs in version {0} cannot be empty.".format(version)

    return None

def no_duplicate_versions_in_changes(changes):
    seen = {}
    for change in changes:
        version = pretty_version(change["version"])
        if version in seen:
            return "Duplicate versions {0}.".format(version)
        seen[version] = True

    return None

def assert_latest_ordered_versions(changes):
    versions = [change["version"] for change in changes[::-1]]
    for i in range(1, len(versions)):
        previous = versions[i - 1]
        current = versions[i]
        is_invalid = False
        
        if previous["major"] > current["major"]:
            is_invalid = True
        elif previous["major"] == current["major"] and previous["minor"] > current["minor"]:
            is_invalid = True

        if is_invalid:
            return "Version {0} must come before {1}.".format(pretty_version(previous), pretty_version(current))

    return None

# Schemas ----------------------------------------------------------------------

Version = Schema(
    fields = {
        "major": Int(required = True),
        "minor": Int(required = True),
    }
)

Change = Schema(
    fields = {
        "version": Object(Version, required = True),
        "log": String(),
    },
    validations = [no_empty_log_in_change]
)

ChangeLog = Schema(
    fields = {
        "changes": List(Change, validations = [
            no_duplicate_versions_in_changes,
            assert_latest_ordered_versions,
        ])
    }
)
