#!/bin/bash
# -----------------------------------------------------------------------------
# EatSavvy Kubernetes Deployment Script
# 
# SAFETY: This script will ONLY deploy to contexts containing "eatsavvy"
# -----------------------------------------------------------------------------

set -e

# Configuration
KUBECONFIG_FILE="${KUBECONFIG:-$HOME/.kube/eatsavvy.config}"
NAMESPACE="eatsavvy"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CURRENT_CONTEXT=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# -----------------------------------------------------------------------------
# Safety Check Function
# -----------------------------------------------------------------------------
safety_check() {
    echo "=============================================="
    echo "EATSAVVY DEPLOYMENT SAFETY CHECK"
    echo "=============================================="
    
    # Check if kubeconfig file exists
    if [[ ! -f "$KUBECONFIG_FILE" ]]; then
        echo -e "${RED}ERROR: Kubeconfig file not found: $KUBECONFIG_FILE${NC}"
        echo "Run this first:"
        echo "  oci ce cluster create-kubeconfig --cluster-id <id> --file $KUBECONFIG_FILE --region <region> --token-version 2.0.0"
        exit 1
    fi
    
    # Get current context
    CURRENT_CONTEXT=$(KUBECONFIG="$KUBECONFIG_FILE" kubectl config current-context 2>/dev/null)
    
    if [[ -z "$CURRENT_CONTEXT" ]]; then
        echo -e "${RED}ERROR: No current context set in $KUBECONFIG_FILE${NC}"
        exit 1
    fi
    
    echo "Kubeconfig: $KUBECONFIG_FILE"
    echo "Context:    $CURRENT_CONTEXT"
    echo ""
    
    # SAFETY CHECK: Context must contain "eatsavvy"
    if [[ "$CURRENT_CONTEXT" != *"eatsavvy"* ]]; then
        echo -e "${RED}=============================================="
        echo "SAFETY CHECK FAILED!"
        echo "=============================================="
        echo "Context '$CURRENT_CONTEXT' does not contain 'eatsavvy'"
        echo ""
        echo "This prevents accidental deployment to wrong clusters."
        echo ""
        echo "To fix, rename your context:"
        echo "  KUBECONFIG=$KUBECONFIG_FILE kubectl config rename-context $CURRENT_CONTEXT eatsavvy-cluster"
        echo -e "==============================================${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Safety check passed!${NC}"
    echo ""
    
    # Verify cluster connectivity
    echo "Verifying cluster connectivity..."
    if ! KUBECONFIG="$KUBECONFIG_FILE" kubectl cluster-info &>/dev/null; then
        echo -e "${RED}ERROR: Cannot connect to cluster${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Cluster is reachable${NC}"
    echo ""
}

# -----------------------------------------------------------------------------
# Confirmation Prompt
# -----------------------------------------------------------------------------
confirm() {
    local action=$1
    local context=$2
    
    echo ""
    echo -e "${BOLD}=============================================="
    echo "CONFIRMATION REQUIRED"
    echo "==============================================${NC}"
    echo ""
    echo -e "Action:  ${YELLOW}${action}${NC}"
    echo -e "Context: ${YELLOW}${context}${NC}"
    echo ""
    
    read -p "Are you sure you want to continue? [y/N] " -n 1 -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}Aborted.${NC}"
        exit 1
    fi
    echo ""
}

# -----------------------------------------------------------------------------
# Check secrets file exists
# -----------------------------------------------------------------------------
check_secrets() {
    if [[ ! -f "$SCRIPT_DIR/01-secrets.yaml" ]]; then
        echo -e "${RED}ERROR: 01-secrets.yaml not found${NC}"
        echo ""
        echo "Create it from the example:"
        echo "  cp 01-secrets.yaml.example 01-secrets.yaml"
        echo "  # Edit 01-secrets.yaml with your credentials"
        exit 1
    fi
    echo -e "${GREEN}✓ Secrets file found${NC}"
}

# -----------------------------------------------------------------------------
# Deploy Function
# -----------------------------------------------------------------------------
deploy() {
    local manifest=$1
    echo -e "${YELLOW}Applying: $manifest${NC}"
    KUBECONFIG="$KUBECONFIG_FILE" kubectl apply -f "$manifest"
}

# -----------------------------------------------------------------------------
# Main
# -----------------------------------------------------------------------------
main() {
    local action="${1:-check}"
    
    case "$action" in
        apply|deploy)
            safety_check
            check_secrets
            
            confirm "DEPLOY resources to cluster" "$CURRENT_CONTEXT"
            
            echo "=============================================="
            echo "DEPLOYING EATSAVVY TO KUBERNETES"
            echo "=============================================="
            echo -e "Context: ${YELLOW}${CURRENT_CONTEXT}${NC}"
            echo ""
            
            # Apply in order (dependencies first)
            deploy "$SCRIPT_DIR/00-namespace.yaml"
            deploy "$SCRIPT_DIR/01-secrets.yaml"
            deploy "$SCRIPT_DIR/02-rabbitmq.yaml"
            deploy "$SCRIPT_DIR/03-api.yaml"
            deploy "$SCRIPT_DIR/04-worker.yaml"
            deploy "$SCRIPT_DIR/05-cloudflared.yaml"
            
            echo ""
            echo -e "${GREEN}=============================================="
            echo "DEPLOYMENT COMPLETE!"
            echo "==============================================${NC}"
            echo ""
            echo "Check status:"
            echo "  KUBECONFIG=$KUBECONFIG_FILE kubectl get pods -n $NAMESPACE"
            ;;
            
        delete|destroy)
            safety_check
            
            confirm "DELETE all resources from cluster" "$CURRENT_CONTEXT"
            
            echo "=============================================="
            echo "DELETING EATSAVVY FROM KUBERNETES"
            echo "=============================================="
            echo -e "Context: ${YELLOW}${CURRENT_CONTEXT}${NC}"
            echo ""
            
            # Delete in reverse order
            KUBECONFIG="$KUBECONFIG_FILE" kubectl delete -f "$SCRIPT_DIR/05-cloudflared.yaml" --ignore-not-found
            KUBECONFIG="$KUBECONFIG_FILE" kubectl delete -f "$SCRIPT_DIR/04-worker.yaml" --ignore-not-found
            KUBECONFIG="$KUBECONFIG_FILE" kubectl delete -f "$SCRIPT_DIR/03-api.yaml" --ignore-not-found
            KUBECONFIG="$KUBECONFIG_FILE" kubectl delete -f "$SCRIPT_DIR/02-rabbitmq.yaml" --ignore-not-found
            KUBECONFIG="$KUBECONFIG_FILE" kubectl delete -f "$SCRIPT_DIR/01-secrets.yaml" --ignore-not-found
            KUBECONFIG="$KUBECONFIG_FILE" kubectl delete -f "$SCRIPT_DIR/00-namespace.yaml" --ignore-not-found
            
            echo -e "${GREEN}Deletion complete!${NC}"
            ;;
            
        status)
            safety_check
            echo -e "Context: ${YELLOW}${CURRENT_CONTEXT}${NC}"
            echo ""
            KUBECONFIG="$KUBECONFIG_FILE" kubectl get all -n $NAMESPACE
            ;;
            
        logs)
            local component="${2:-api}"
            safety_check
            echo -e "Context: ${YELLOW}${CURRENT_CONTEXT}${NC}"
            echo ""
            KUBECONFIG="$KUBECONFIG_FILE" kubectl logs -n $NAMESPACE -l app="$component" -f
            ;;
            
        check)
            safety_check
            check_secrets
            echo -e "${GREEN}All checks passed! Ready to deploy.${NC}"
            ;;
            
        *)
            echo "Usage: $0 {apply|delete|status|logs [component]|check}"
            echo ""
            echo "Commands:"
            echo "  apply   - Deploy all manifests"
            echo "  delete  - Delete all resources"
            echo "  status  - Show deployment status"
            echo "  logs    - Tail logs (default: api, or specify: rabbitmq, worker, cloudflared)"
            echo "  check   - Run safety check only (default)"
            echo ""
            echo "Environment:"
            echo "  KUBECONFIG - Path to kubeconfig (default: ~/.kube/eatsavvy.config)"
            exit 1
            ;;
    esac
}

main "$@"
