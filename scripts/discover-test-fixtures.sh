#!/bin/bash
set -e

# Check for required environment variable
if [ -z "$SIMPLEMDM_APIKEY" ]; then
    echo "❌ Error: SIMPLEMDM_APIKEY environment variable not set"
    echo ""
    echo "Usage:"
    echo "  export SIMPLEMDM_APIKEY='your-api-key'"
    echo "  ./scripts/discover-test-fixtures.sh"
    exit 1
fi

echo "🔍 Discovering SimpleMDM test fixtures..."
echo ""

# Function to make API calls
api_call() {
    curl -s -u "${SIMPLEMDM_APIKEY}:" "https://a.simplemdm.com/api/v1/$1"
}

# 1. Find enrolled device
echo "1️⃣  Finding enrolled device..."
DEVICE_RESPONSE=$(api_call "devices")
DEVICE_ID=$(echo "$DEVICE_RESPONSE" | jq -r '.data[0].id // empty')

if [ -n "$DEVICE_ID" ]; then
    DEVICE_NAME=$(echo "$DEVICE_RESPONSE" | jq -r '.data[0].attributes.name // "Unknown"')
    echo "   ✅ Found device: $DEVICE_NAME (ID: $DEVICE_ID)"
else
    echo "   ❌ No enrolled devices found"
    echo "   ℹ️  You need at least one enrolled device in SimpleMDM"
fi
echo ""

# 2. Find device group
echo "2️⃣  Finding device group for cloning..."
GROUP_RESPONSE=$(api_call "device_groups")
GROUP_ID=$(echo "$GROUP_RESPONSE" | jq -r '.data[0].id // empty')

if [ -n "$GROUP_ID" ]; then
    GROUP_NAME=$(echo "$GROUP_RESPONSE" | jq -r '.data[0].attributes.name // "Unknown"')
    echo "   ✅ Found device group: $GROUP_NAME (ID: $GROUP_ID)"
else
    echo "   ❌ No device groups found"
    echo "   ℹ️  Device groups are created automatically, but you can create custom ones"
fi
echo ""

# 3. Find script job
echo "3️⃣  Finding script job..."
SCRIPT_JOB_RESPONSE=$(api_call "script_jobs")
SCRIPT_JOB_ID=$(echo "$SCRIPT_JOB_RESPONSE" | jq -r '.data[0].id // empty')

if [ -n "$SCRIPT_JOB_ID" ]; then
    SCRIPT_JOB_STATUS=$(echo "$SCRIPT_JOB_RESPONSE" | jq -r '.data[0].attributes.status // "Unknown"')
    echo "   ✅ Found script job: Status $SCRIPT_JOB_STATUS (ID: $SCRIPT_JOB_ID)"
else
    echo "   ⚠️  No script jobs found"
    echo "   ℹ️  To create one: Run any script via SimpleMDM UI on a device/group"
fi
echo ""

# 4. Find DDM-capable device
echo "4️⃣  Finding DDM-capable device (macOS 13+, iOS 15+)..."
# Use same device as #1 for simplicity (most modern devices support DDM)
if [ -n "$DEVICE_ID" ]; then
    DEVICE_OS=$(echo "$DEVICE_RESPONSE" | jq -r '.data[0].attributes.os_version // "Unknown"')
    echo "   ✅ Using device: $DEVICE_NAME (OS: $DEVICE_OS, ID: $DEVICE_ID)"
    echo "   ℹ️  Ensure OS version supports DDM: macOS 13+, iOS 15+"
else
    echo "   ❌ No devices available for DDM testing"
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 SUMMARY - GitHub Secrets Commands"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ -n "$DEVICE_ID" ] && [ -n "$GROUP_ID" ] && [ -n "$SCRIPT_JOB_ID" ]; then
    echo "✅ All required fixture IDs found!"
    echo ""
    echo "Run these commands to set GitHub secrets:"
    echo ""
    echo "gh secret set SIMPLEMDM_DEVICE_ID --body \"$DEVICE_ID\""
    echo "gh secret set SIMPLEMDM_DEVICE_GROUP_CLONE_SOURCE_ID --body \"$GROUP_ID\""
    echo "gh secret set SIMPLEMDM_SCRIPT_JOB_ID --body \"$SCRIPT_JOB_ID\""
    echo "gh secret set SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID --body \"$DEVICE_ID\""
    echo ""
    echo "Or copy/paste this single command:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    cat << EOF
gh secret set SIMPLEMDM_DEVICE_ID --body "$DEVICE_ID" && \\
gh secret set SIMPLEMDM_DEVICE_GROUP_CLONE_SOURCE_ID --body "$GROUP_ID" && \\
gh secret set SIMPLEMDM_SCRIPT_JOB_ID --body "$SCRIPT_JOB_ID" && \\
gh secret set SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID --body "$DEVICE_ID"
EOF
    echo ""
else
    echo "⚠️  Missing some fixtures. Please review the output above."
    echo ""
    echo "Available commands:"
    [ -n "$DEVICE_ID" ] && echo "  gh secret set SIMPLEMDM_DEVICE_ID --body \"$DEVICE_ID\""
    [ -n "$GROUP_ID" ] && echo "  gh secret set SIMPLEMDM_DEVICE_GROUP_CLONE_SOURCE_ID --body \"$GROUP_ID\""
    [ -n "$SCRIPT_JOB_ID" ] && echo "  gh secret set SIMPLEMDM_SCRIPT_JOB_ID --body \"$SCRIPT_JOB_ID\""
    [ -n "$DEVICE_ID" ] && echo "  gh secret set SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID --body \"$DEVICE_ID\""
    echo ""
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📖 Next Steps:"
echo "   1. Run the gh secret set commands above"
echo "   2. Verify secrets: gh secret list"
echo "   3. Re-run GitHub Actions workflow"
echo "   4. See TESTING_SETUP.md for detailed documentation"
echo ""
echo "✨ Done!"