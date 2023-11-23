class QuoteDataRepo {
    constructor(title="", commission="") {
        this.casingList = [];
        this.fixtureList = [];
        this.productsCollection = {};
        this.commission = commission;
        this.title = title;
    }

    addCasingProduct(casingObj) {
        this.casingList.push(casingObj);
        this.productsCollection[casingObj.uuid] = casingObj;
    }

    addFixtureProduct(fixtureObj){
        this.fixtureList.push(fixtureObj);
        this.productsCollection[fixtureObj.uuid] = fixtureObj;
    }

    setupCasingSelectors() {
        for(let fixture of this.fixtureList) {
            fixture.setupCasingSelector(this.casingList, this.productsCollection);
        }
    }

    renderAllProducts(parent) {
        for (const [uuid, productItem] of Object.entries(this.productsCollection)) {
            productItem.render(parent);
        }
    }
}