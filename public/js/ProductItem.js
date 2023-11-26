class ProductItem {
    constructor(uuid, prodData, position) {
        this.uuid = uuid;
        this.position = position;
        this.prodData = prodData;
        this.node = null;
        this.group = "A";  // default group
    }

    //custom event system
    on(name, callback) {
        var callbacks = this[name];
        if (!callbacks) this[name] = [callback];
        else callbacks.push(callback);
    }

    dispatch(event) {
        var callbacks = this[event.name];
        if (callbacks) callbacks.forEach(callback => callback(event));
    }

    setupNode() {
        let templateNode = document.querySelector("#prod-item-template")
            .content
            .querySelector(".prod-list-item"); //get the li element

        this.node = templateNode.cloneNode(true);

        //setup event handlers
        let self = this;
        
        this.node.querySelector(".prod-group").addEventListener(
            "change", function () {
                self.handleGroupSelectChange.call(self);
            }
        );

    }

    handleGroupSelectChange() {
        //this.group = this.node.querySelector(".prod-group").value = "B";
    }

    render(parent) {
        this.node.querySelector(".width").textContent = this.prodData.width;
        this.node.querySelector(".height").textContent = this.prodData.height;
        this.node.querySelector(".prod-id").textContent = this.prodData.product_id;

        parent.appendChild(this.node);
    }

    toJson() {
        return {
            "product_data": this.prodData,
            "position": parseInt(this.position),
            "group": this.group,
        }
    }

}