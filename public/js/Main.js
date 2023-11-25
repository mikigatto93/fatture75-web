const input = document.getElementById("data-file");
const submitBtn = document.getElementById("submit-1");
const quoteDataContainer = document.getElementById("quote-data-container");
const quoteHeader = document.getElementById("quote-header");
const quoteItemsList = document.getElementById("quote-items-list");
const sendDataBtn = document.getElementById("send-data-btn");

let QuoteData = new QuoteDataRepo();


submitBtn.addEventListener("click", () => {
  let formData = new FormData();
  formData.append("file", input.files[0]);

  let promise = makeProcessQuoteRequest(formData);

  promise.then(
    (data) => { processQuoteData(data) }
  );

});


async function makeProcessQuoteRequest(formData) {
  let response = await fetch(
    "http://localhost:8888/api/quote_data",
    {
      method: "POST",
      body: formData,
    }
  );

  return await response.json();
}


function processQuoteData(data) {
  console.log(data);
  quoteHeader.textContent = data.title;

  //format header
  QuoteData.commission = data.commission;
  QuoteData.title = data.title;

  //format products
  for (let position in data.products) {
    let prod = data.products[position];
    let prodData = formatProduct(prod);
    createProductItem(prodData, position);
  }

  QuoteData.setupCasingSelectors();

  QuoteData.renderAllProducts(quoteItemsList);
}

function formatProduct(product) {
  let prodObj = {
    "height": product.Height,
    "width": product.Width,
    "notes": product.Notes,
    "tot_price": product.Price,
    "product_id": product.ProductId,
    "quantity": product.Quantity,
    "reference": product.Reference,
  };

  if (product.Depth > 0) {
    prodObj["depth"] = parseInt(product.Depth);
  }
  
  return prodObj;
}

function createProductItem(product, position) {
  let prodElement;
  const uuid = crypto.randomUUID();

  if ("depth" in product) {
    prodElement = new CasingProductItem(uuid, product, position);
    QuoteData.addCasingProduct(prodElement);
  } else {
    prodElement = new FixtureProductItem(uuid, product, position);
    QuoteData.addFixtureProduct(prodElement);
  }
  
  prodElement.setupNode();
}
